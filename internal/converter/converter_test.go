package converter_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nhatthm/moneylovercli-plugin-n26/internal/converter"
	"github.com/stretchr/testify/assert"
)

type delayedReader struct {
	*strings.Reader
	wait time.Duration
}

func (r *delayedReader) Read(p []byte) (n int, err error) {
	<-time.After(r.wait)

	return r.Reader.Read(p)
}

func newDelayedReader(wait time.Duration, s string) *delayedReader {
	return &delayedReader{
		Reader: strings.NewReader(s),
		wait:   wait,
	}
}

type readError struct {
	err error
}

func (r *readError) Read([]byte) (n int, err error) {
	return 0, r.err
}

func newReaderHasError(err error) *readError {
	return &readError{err: err}
}

type writeError struct {
	err error
}

func (r *writeError) Write([]byte) (n int, err error) {
	return 0, r.err
}

func newWriterHasError(err error) *writeError {
	return &writeError{err: err}
}

type buffer struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

func (b *buffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buffer.Write(p)
}

func (b *buffer) Bytes() []byte {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buffer.Bytes()
}

func TestConvert(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		context        func() (context.Context, context.CancelFunc)
		input          io.Reader
		output         io.Writer
		expectedOutput []byte
		expectedError  string
	}{
		{
			scenario:      "input is nil",
			input:         bytes.NewReader(nil),
			expectedError: `could not get json token: EOF`,
		},
		{
			scenario:      "input is not a json",
			input:         strings.NewReader("hello"),
			expectedError: `could not get json token: invalid character 'h' looking for beginning of value`,
		},
		{
			scenario:      "input is an error string",
			input:         strings.NewReader(`could not get transactions: could not login"`),
			expectedError: `could not get json token: invalid character 'c' looking for beginning of value`,
		},
		{
			scenario:      "input is a json string",
			input:         strings.NewReader(`"hello"`),
			expectedError: "input is not an array",
		},
		{
			scenario:      "input is an integer",
			input:         strings.NewReader(`42`),
			expectedError: "input is not an array",
		},
		{
			scenario:      "input is a boolean",
			input:         strings.NewReader(`true`),
			expectedError: "input is not an array",
		},
		{
			scenario:      "input is an object",
			input:         strings.NewReader(`{}`),
			expectedError: "input is not an array",
		},
		{
			scenario: "input has empty space prefix",
			input:    strings.NewReader(`  []`),
		},
		{
			scenario: "context is canceled",
			context: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				return ctx, func() {}
			},
			input:         strings.NewReader(`[]`),
			expectedError: `decode interrupted`,
		},
		{
			scenario: "read timed out",
			context: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Millisecond)
			},
			input:         newDelayedReader(15*time.Millisecond, `[{},{}]`),
			expectedError: `decode interrupted`,
		},
		{
			scenario:      "read error",
			input:         newReaderHasError(errors.New(`read error`)),
			expectedError: `could not get json token: read error`,
		},
		{
			scenario:      "write error",
			input:         newDelayedReader(5*time.Millisecond, `[{},{},{},{}]`),
			output:        newWriterHasError(errors.New(`write error`)),
			expectedError: `could not write transaction: write error`,
		},
		{
			scenario:      "input has an invalid element",
			input:         strings.NewReader(`  [true]`),
			expectedError: `could not decode transaction: json: cannot unmarshal bool into Go value of type transaction.Transaction`,
		},
		{
			scenario: "input is valid",
			input: strings.NewReader(`
				[
					{
						"id": "801d35f4-f550-446a-974a-0d5dc2c1f55d",
						"userId": "7e3f710b-349d-4203-9c5d-cfbc716e1b8e",
						"type": "CT",
						"amount": 10,
						"currencyCode": "EUR",
						"visibleTS": 1617631557000,
						"partnerBic": "NTSBDEB1XXX",
						"partnerName": "Jane Doe",
						"accountId": "98f0afa3-e906-493a-a37f-afe29c7f9f2e",
						"partnerIban": "DEXX1001100126XXXXXXXX",
						"partnerBankName": "Revolut",
						"category": "micro-v2-income",
						"cardId": "f2252c42-c188-4b43-ab68-131024782b3d",
						"referenceText": "A random transaction",
						"userCertified": 1617545157000,
						"pending": false,
						"transactionNature": "NORMAL",
						"createdTS": 1617541557000,
						"smartLinkId": "fcdec3cb-47b2-4ca3-b98d-b326e1cc5a0c",
						"smartContactId": "3edce485-6853-40bf-aa08-309c2eb3e7db",
						"linkId": "6f06f5fb-074d-4242-b280-db2af2fe6405",
						"confirmed": 1617545157000
					},
					{
						"id": "b7139067-6a3c-42ec-9c91-2c3dfdef2ece",
						"userId": "7e3f710b-349d-4203-9c5d-cfbc716e1b8e",
						"type": "PT",
						"amount": -15.6,
						"currencyCode": "EUR",
						"originalAmount": -15.6,
						"originalCurrency": "EUR",
						"exchangeRate": 1,
						"merchantCity": "GBR",
						"visibleTS": 1618494602996,
						"mcc": 4829,
						"mccGroup": 10,
						"merchantName": "A Random Merchant",
						"accountId": "98f0afa3-e906-493a-a37f-afe29c7f9f2e",
						"category": "micro-v2-insurances-finances",
						"cardId": "5bf25fb1-a578-43cb-bc32-47e7443ea46d",
						"userCertified": 1618556043333,
						"pending": false,
						"transactionNature": "NORMAL",
						"createdTS": 1618556043337,
						"merchantCountry": 28,
						"merchantCountryCode": 440,
						"smartLinkId": "ff756c53-6ece-4930-91c1-496aa9c0b364",
						"linkId": "bbc439cc-d074-44cb-af44-2873e96c8f07",
						"txnCondition": "ECOMMERCE",
						"confirmed": 1618556043333
					}
				]
			`),
			expectedOutput: []byte(`{
    "id": "801d35f4-f550-446a-974a-0d5dc2c1f55d",
    "accountName": "",
    "accountId": "98f0afa3-e906-493a-a37f-afe29c7f9f2e",
    "accountBank": "N26",
    "amount": 10,
    "category": "micro-v2-income",
    "referenceText": "A random transaction",
    "displayDate": "2021-04-05T14:05:57Z",
    "partnerName": "Jane Doe",
    "partnerId": "DEXX1001100126XXXXXXXX",
    "partnerBank": "Revolut"
}
{
    "id": "b7139067-6a3c-42ec-9c91-2c3dfdef2ece",
    "accountName": "",
    "accountId": "98f0afa3-e906-493a-a37f-afe29c7f9f2e",
    "accountBank": "N26",
    "amount": -15.6,
    "category": "micro-v2-insurances-finances",
    "referenceText": "",
    "displayDate": "2021-04-15T13:50:02.996Z",
    "partnerName": "A Random Merchant",
    "partnerId": "",
    "partnerBank": ""
}
`),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			if tc.context == nil {
				tc.context = func() (context.Context, context.CancelFunc) {
					return context.Background(), func() {}
				}
			}

			buf := new(buffer)

			if tc.output == nil {
				tc.output = buf
			} else {
				tc.output = io.MultiWriter(tc.output, buf)
			}

			ctx, cancel := tc.context()
			defer cancel()

			err := converter.Convert(ctx, tc.input, tc.output, converter.WithPretty(true))
			result := buf.Bytes()

			if tc.expectedOutput == nil {
				assert.Nil(t, result)
			} else {
				t.Log(string(result))

				assert.NotNil(t, result)
				assert.Equal(t, string(tc.expectedOutput), string(result))
			}

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
