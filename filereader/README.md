# filereader

The filereader is the main routine for parsing a given file, and returning the requested results.  It includes unit tests (in `filereader_test.go`) which take a sampling of known data, and verify the correct responses.

# Operation

The filereader starts from the end of the file and reads backwards in chunks (of configurable size).  The filereader works back through the file until the beginning of the file is reached, or the requested number of matches have been found.

## Assumptions

The filereader assumes that if no matching text is supplied, then all lines are a match.  

The filereader assumes that if a `count` argument of 0 (or negative) is provided, then all matches should be returned.  A positive integer will set the limit on the number of matches returned.

## Design considerations

The approach to read from the end of the file was driven by the following requirements:
- the lines returned must be presented with the newest log events first
- log files will be written with the newest events at the end of the file
- the service should work and be reasonbly performant when requesting files of >1GB

This combination of requirements led to the decision of starting at the end of the file and working backwards.  Starting at the beginning of the file and working forwards has the following consequences

- The entire file must always be read. While this could still be true working backwards, in the case where a count argument is provided, then the entire file does not have to be read if the count is met prior to reaching the beginning.  In the case of a large file with a small count, then this can and will result in a significant savings.
- The entire file might need to be cached.  In the scenario with a very broad (or null) match string, and a large count, in order to reverse  the file line-by-line, the entire file would need to be cached.  This would either be in memory, which would result in very large memory consumption (for large file), or in a temporary file which would require rewriting the entire file in reverse order, and then re-reading to return the file.

## Testing considerations

The unit tests contain a sample file with known input.  The input is supplied to the file reading code, and the return values are checked.
The test cases include
- That the entire file can be read in a single read (small file read)
- That the entire file is properly read when the read size is smaller than a line length (short reads)
- The file is properly parsed when the read chunks fall on line boundaries
- The count parameter is properly honored
- the match parameter is properly honored
- 