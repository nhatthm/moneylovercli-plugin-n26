package converter

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/bool64/ctxd"
	mlapi "github.com/nhatthm/moneyloverapi/pkg/transaction"
	n26api "github.com/nhatthm/n26api/pkg/transaction"
)

// Option is option to setup the converter.
type Option func(c *Converter)

// Converter is converter option.
type Converter struct {
	Encoder *json.Encoder
	Decoder *json.Decoder
}

// Convert converts n26 transactions from input stream and write the result to output stream.
func Convert(ctx context.Context, r io.Reader, w io.Writer, opts ...Option) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c := newConverter(r, w, opts...)

	resCh, readErr := read(ctx, c.Decoder)
	writeErr := write(ctx, c.Encoder, resCh)

	running := map[string]bool{"read": true, "write": true}

	for {
		select {
		case err, ok := <-readErr:
			if ok {
				return err
			}

			delete(running, "read")

		case err, ok := <-writeErr:
			if ok {
				return err
			}

			delete(running, "write")

		default:
		}

		if len(running) == 0 {
			break
		}
	}

	return nil
}

// WithPretty sets pretty json output.
func WithPretty(pretty bool) Option {
	return func(c *Converter) {
		if pretty {
			c.Encoder.SetIndent("", "    ")
		}
	}
}

func newConverter(r io.Reader, w io.Writer, opts ...Option) *Converter {
	c := &Converter{
		Encoder: json.NewEncoder(w),
		Decoder: json.NewDecoder(r),
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

func read(ctx context.Context, dec *json.Decoder) (<-chan n26api.Transaction, <-chan error) {
	resCh := make(chan n26api.Transaction)
	errCh := make(chan error, 1)

	// nolint: errcheck
	go func() (err error) {
		defer close(resCh)
		defer func() {
			if err != nil {
				errCh <- err
			}

			close(errCh)
		}()

		token, err := dec.Token()
		if err != nil {
			return ctxd.WrapError(ctx, err, "could not get json token")
		}

		if token != json.Delim('[') {
			return ctxd.NewError(ctx, "input is not an array", "token", token)
		}

		idx := 0

		for {
			select {
			case <-ctx.Done():
				return ctxd.NewError(ctx, "decode interrupted")

			default:
				if !dec.More() {
					return nil
				}

				var t n26api.Transaction

				if err := dec.Decode(&t); err != nil {
					return ctxd.WrapError(ctx, err, "could not decode transaction", "index", idx)
				}

				resCh <- t

				idx++
			}
		}
	}()

	return resCh, errCh
}

func write(ctx context.Context, enc *json.Encoder, transactions <-chan n26api.Transaction) <-chan error {
	errCh := make(chan error, 1)

	// nolint: errcheck
	go func() (err error) {
		defer func() {
			if err != nil {
				errCh <- err
			}

			close(errCh)
		}()

		for {
			select {
			case <-ctx.Done():
				return nil

			case t, ok := <-transactions:
				if !ok {
					return nil
				}

				if err := enc.Encode(convert(t)); err != nil {
					return ctxd.WrapError(ctx, err, "could not write transaction")
				}
			}
		}
	}()

	return errCh
}

func convert(t n26api.Transaction) mlapi.BankTransaction {
	bt := mlapi.BankTransaction{
		ID:            t.ID.String(),
		AccountID:     t.AccountID.String(),
		AccountBank:   "N26",
		Amount:        t.Amount,
		Category:      t.Category,
		ReferenceText: t.ReferenceText,
		DisplayDate:   time.Unix(0, t.VisibleTS*1000000).UTC(), // VisibleTS is in ms, convert to nanoseconds.
	}

	if t.MerchantName != "" {
		bt.PartnerName = t.MerchantName
	} else {
		bt.PartnerName = t.PartnerName
		bt.PartnerID = t.PartnerIban
		bt.PartnerBank = t.PartnerBankName
	}

	return bt
}
