package command_test

import (
	"bytes"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/nhatthm/moneylovercli-plugin-n26/internal/command"
	"github.com/stretchr/testify/assert"
)

type buffer struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

func (b *buffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buffer.Write(p)
}

func (b *buffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buffer.String()
}

func TestNewConvert(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		stdin    io.Reader
		stdout   io.Writer
		args     []string
		expected string
	}{
		{
			scenario: "no pretty",
			stdin: strings.NewReader(`
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
					}
				]
			`),
			expected: `{"id":"801d35f4-f550-446a-974a-0d5dc2c1f55d","accountName":"","accountId":"98f0afa3-e906-493a-a37f-afe29c7f9f2e","accountBank":"N26","amount":10,"category":"micro-v2-income","referenceText":"A random transaction","displayDate":"2021-04-05T14:05:57Z","partnerName":"Jane Doe","partnerId":"DEXX1001100126XXXXXXXX","partnerBank":"Revolut"}
`,
		},
		{
			scenario: "pretty",
			stdin: strings.NewReader(`
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
					}
				]
			`),
			args: []string{"-p"},
			expected: `{
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
`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			stdout := new(buffer)
			cmd := command.New(
				command.WithStdin(tc.stdin),
				command.WithStdout(stdout),
			)

			cmd.SetArgs(append([]string{"convert"}, tc.args...))

			err := cmd.Execute()
			result := stdout.String()

			t.Log(result)

			assert.Equal(t, tc.expected, result)
			assert.NoError(t, err)
		})
	}
}
