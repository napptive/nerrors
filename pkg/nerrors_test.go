package nerrors

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func TestNerror(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Shopping Cart Suite")
}

// function to test if ExtendedError satisfy the error interface
func test() error {
	return NewInternalError("internal error")
}

var _ = ginkgo.Describe("Handler test on nerrors calls", func() {
	ginkgo.Context("Check error lib", func() {
		ginkgo.It("Check the error has StackTrace", func() {
			msg := "unable to find the record"
			err := NewNotFoundError(msg)
			gomega.Expect(err).ShouldNot(gomega.BeNil())
			gomega.Expect(err.StackTrace).ShouldNot(gomega.BeEmpty())
			gomega.Expect(err.Msg).Should(gomega.Equal(msg))
			gomega.Expect(err.Code).Should(gomega.Equal(NotFound))
		})

		ginkgo.It("Check a complex error is well created", func() {
			msg := "unable to find the record"
			err := NewNotFoundError(msg)
			complexErr := NewAbortedErrorFrom(err, "unable to continue")
			gomega.Expect(complexErr).ShouldNot(gomega.BeNil())
			gomega.Expect(complexErr.StackTrace).ShouldNot(gomega.BeEmpty())
			fmt.Println(complexErr.String())
			gomega.Expect(complexErr.String()).Should(gomega.ContainSubstring(msg))
			gomega.Expect(complexErr.Code).Should(gomega.Equal(Aborted))
			gomega.Expect(complexErr.From).ShouldNot(gomega.BeNil())
			gomega.Expect(complexErr.From).Should(gomega.Equal(err))
		})

		ginkgo.It("Check if ExtendedError is an error", func() {
			var err error
			err = NewInternalError("message")
			gomega.Expect(err).NotTo(gomega.BeNil())

			errCheck := test()
			gomega.Expect(errCheck).NotTo(gomega.BeNil())
			gomega.Expect(errCheck.Error()).NotTo(gomega.BeEmpty())
		})
	})
	ginkgo.Context("checking conversions to GRPC", func() {
		ginkgo.It("can convert to GRPC error", func(){
			err := NewInternalError("internal error")
			grpcError := err.ToGRPC()
			gomega.Expect(grpcError).NotTo(gomega.BeNil())
		})
		ginkgo.It("can convert a complex error to GRPC error", func(){
			common := status.Errorf(codes.Aborted, "new error")
			err := NewInternalErrorFrom(common,"internal error")
			grpcError := err.ToGRPC()
			gomega.Expect(grpcError).NotTo(gomega.BeNil())
		})
	})
})

// ToGrpc FromGrpc

