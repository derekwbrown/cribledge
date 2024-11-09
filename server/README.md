# server

The server module implements a simple HTTP server.  It's sole function is to provide the HTTP endpoint for answering the REST requests, and parse the various arguments.


## Testing considerations

The server unit tests provide a more complete set of testing scenarios with different data files.  The data files were taken from a Linux device, and are the `dmesg` and `syslog` files.  These files provide interesting data because in addition to having many similar messages, the entries are timestamped which provides an easy way to visually verify that the log lines are presented in reverse order.

The unit tests contain a sample file with known input.  The input is supplied to the file reading code, and the return values are checked.
The test cases include
- That an entire file is returned in reverse order
- That the correct number of matches is returned.  
- That the correct number of matches is limited by the `count` parameter
- That the server properly reads files with path separators, and reads into subdirectories.

The test does set up a complete http server, so if the test environment is inadequate (such as port 8080 already being in use), the test will fail.
