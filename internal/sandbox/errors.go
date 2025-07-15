package sandbox

import (
	"fmt"

	"github.com/sourcenetwork/acp_core/pkg/errors"
)

func itoa(i uint64) string {
	return fmt.Sprintf("%d", i)
}

func newNewSandboxErr(err error) error {
	return errors.Wrap("NewSandbox failed", err)
}

func newListSandboxesErr(err error) error {
	return errors.Wrap("ListSandboxes failed", err)
}

func newSetStateErr(err error, handle uint64) error {
	return errors.Wrap("SetState failed", err, errors.Pair("handle", itoa(handle)))
}

func newVerifyTheoremsErr(err error, handle uint64) error {
	return errors.Wrap("VerifyTheorems failed", err, errors.Pair("handle", itoa(handle)))
}

func newGetCatalogueErr(err error, handle uint64) error {
	return errors.Wrap("GetCatalogue failed", err, errors.Pair("handle", itoa(handle)))
}

func newRestoreScratchpadErr(err error, handle uint64) error {
	return errors.Wrap("RestoreScratchpad failed", err, errors.Pair("handle", itoa(handle)))
}

func newGetSandboxErr(err error, handle uint64) error {
	return errors.Wrap("GetSandbox failed", err, errors.Pair("handle", itoa(handle)))
}

func newExpandError(err error, handle uint64) error {
	return errors.Wrap("Expand failed", err, errors.Pair("handle", itoa(handle)))
}
