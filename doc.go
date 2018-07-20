package errors

// A package designed to improve error handling in go.
//
// It revolves around a few simple observations:
//
// Errors always have a root cause, this is the base case such as io.EOF.
// It allows API's to document specific errors so they can be handled.
//
// Errors should have a message for humans, this should be informative.
//
// Errors should optionally have some structured context for debugging.
//
// Errors should have a stack trace from the origin to where it was handled for debugging
// purposes.
