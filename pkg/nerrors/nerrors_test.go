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
	ginkgo.RunSpecs(t, "Nerror Suite")
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
			var err error = NewInternalError("message")
			gomega.Expect(err).NotTo(gomega.BeNil())

			errCheck := test()
			gomega.Expect(errCheck).NotTo(gomega.BeNil())
			gomega.Expect(errCheck.Error()).NotTo(gomega.BeEmpty())
		})

		ginkgo.It("Check the error has StackTrace (FromError conversion)", func() {
			msg := "unable to find the record"
			err :=  FromError(fmt.Errorf(msg))
			gomega.Expect(err).ShouldNot(gomega.BeNil())
			gomega.Expect(err.StackTrace).ShouldNot(gomega.BeEmpty())
		})
	})
	// ToGrpc
	ginkgo.Context("checking conversions to GRPC", func() {
		ginkgo.It("can convert to GRPC error", func() {
			err := NewInternalError("internal error")
			grpcError := err.ToGRPC()
			gomega.Expect(grpcError).NotTo(gomega.BeNil())
		})
		ginkgo.It("can convert a complex error to GRPC error", func() {
			common := status.Errorf(codes.Aborted, "new error")
			err := NewInternalErrorFrom(common, "internal error")
			grpcError := err.ToGRPC()
			gomega.Expect(grpcError).NotTo(gomega.BeNil())
		})
		ginkgo.It("can convert a Extended error to GRPC error and the result to error again", func() {
			err := NewInternalError("internal error")
			grpcError := err.ToGRPC()
			gomega.Expect(grpcError).NotTo(gomega.BeNil())

			converted := FromGRPC(grpcError)
			gomega.Expect(converted).ShouldNot(gomega.BeNil())

			gomega.Expect(err).Should(gomega.Equal(converted))

		})
		ginkgo.It("can convert a complex Extended error to GRPC error and the result to error again", func() {
			common := NewNotFoundError("not found")
			err := NewInternalErrorFrom(common, "internal error")
			grpcError := err.ToGRPC()
			gomega.Expect(grpcError).NotTo(gomega.BeNil())

			converted := FromGRPC(grpcError)
			gomega.Expect(converted).ShouldNot(gomega.BeNil())

			gomega.Expect(err).Should(gomega.Equal(converted))

		})
		ginkgo.It("can convert a complex Extended error to GRPC error and the result to error again", func() {
			// internal - not found - fmt.Error
			err := NewInternalErrorFrom(NewNotFoundErrorFrom(fmt.Errorf("Fmt error"), "not found"), "internal error")
			grpcError := err.ToGRPC()
			gomega.Expect(grpcError).NotTo(gomega.BeNil())

			converted := FromGRPC(grpcError)
			gomega.Expect(converted).ShouldNot(gomega.BeNil())

			gomega.Expect(err.Code).Should(gomega.Equal(converted.Code))

		})
		ginkgo.It("can convert a simple error to grpc and convert into extended error again", func() {
			// internal - not found - fmt.Error
			err := fmt.Errorf("Fmt error")
			grpcError := FromError(err).ToGRPC()
			gomega.Expect(grpcError).NotTo(gomega.BeNil())

			converted := FromGRPC(grpcError)
			gomega.Expect(converted).ShouldNot(gomega.BeNil())

		})
	})
	// FromGrpc
	ginkgo.Context("checking conversions from GRPC", func() {
		ginkgo.It("can convert from GRPC", func() {
			err := status.Error(codes.NotFound, "id was not found")
			extended := FromGRPC(err)
			gomega.Expect(extended.Code).Should(gomega.Equal(NotFound))
			gomega.Expect(extended.Msg).ShouldNot(gomega.BeEmpty())
			gomega.Expect(extended.From).Should(gomega.BeNil())
			// Extended Message ALWAYS have Stack Trace
			gomega.Expect(extended.StackTrace).ShouldNot(gomega.BeNil())
		})
		ginkgo.It("can convert from GRPC and to grpc again", func() {
			err := status.Error(codes.NotFound, "id was not found")
			extended := FromGRPC(err)
			gomega.Expect(extended.Code).Should(gomega.Equal(NotFound))

			grpcError := extended.ToGRPC()
			gomega.Expect(grpcError).ShouldNot(gomega.BeNil())
			gomega.Expect(grpcError.Error()).Should(gomega.Equal(err.Error()))

		})
	})
	// FromError
	ginkgo.Context("checking conversions from standard error", func() {
		ginkgo.It("Extended error can be converted to Extended error", func() {
			extended := test() // test returns an ExtendedError as error
			converted := FromError(extended)
			gomega.Expect(converted).Should(gomega.Equal(extended))
		})
		ginkgo.It("Extended error can be converted to Extended error", func() {
			parent := test() // test returns an ExtendedError as error
			extended := NewNotFoundErrorFrom(parent, "not found error")
			converted := FromError(extended)
			gomega.Expect(converted).Should(gomega.Equal(extended))
		})
		ginkgo.It("error can be converted to Extended error", func() {
			standard := fmt.Errorf("standard error")
			converted := FromError(standard)
			gomega.Expect(converted).ShouldNot(gomega.Equal(standard))
			gomega.Expect(converted.StackTrace).ShouldNot(gomega.BeNil())
			gomega.Expect(converted.From).Should(gomega.BeNil())
		})
	})
})



