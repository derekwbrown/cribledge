# cribledge

Welcome to my take home assignment for the Cribbl Edge team.  

The project assumes the following tools are installed and in the path
- Go 1.22 or higher
- Wix 3.11 (for the installer)
- Gnu make
  - The makefile does assume that, along with Make, you have a unix-like environment (msys,cygwin, etc.) such that commands like `rm -f ` work as expected.  This requiremenet is especially for the `clean` target.


# Project Layout

The project is laid out in the following directory structure:

## Directories
### Root Directory
The root directory contains the main application component.  This is the main executable entry point.  The root directory also contains the Makefile, from which all artifacts can be generated, and the tests can be run.

### filereader
This module is the code that implements the logic to read the requested file, and return the matched line(s) in the correct order.  Specific details for the filereader implementation are provided in [the readme](./filereader/README.md)

### server
This module provides the http server interface for answering the REST API requests.  Specific details are provided in [the readme](./server/)

### installer
This directory contains a minimal Windows installer implemented using WiX.  The installer includes little GUI and no options.  It installs and starts the service under the Windows Service Control Manager.  As a sample, the MSI also includes a text file (romeo.txt, which is the text of "Romeo and Juliet").  This allows for a user to install the project, and then test to see specific responses with a known file.

## Makefile
The Makefile includes the following targets:

### all
The `all` target compiles the target executable (cribledge.exe), and then builds the installer.  

### exe
The `exe` target builds only the executable (cribledge.exe), with no installer.  The exe does not need to be run as a Windows service.  It can be run as a command-line application.

### test
The `test` target executes the built in unit tests for the `server` and `filereader` components.  More details on the tests can be found in the respective READMEs.

### clean
The `clean` target removes all target and intermediate files.


# Installation

The application can be installed by simply running the MSI.  It does require administrative privilege, because it will install the application as a Windows service.

# Known issues

A number of issues are not covered in this sample application because they are out of scope of the assigned test.  Some of these issues include:

## Port assignment

The server listens on a hard-coded port (8080).  If that port is in use by another process, then the service will fail.  There is currently no ability to configure an alternate port.

## Firewall accessibility

The installation does not attempt to open a hole in the firewall to allow the traffic for the default port (8080).  The user therefore, depending on host configuration, may have to manually create a firewall rule to allow the traffic to reach the target.

## Privilege escalation

The server runs as `LOCAL_SYSTEM` (root).  However, the API requests have no authentication requirements.  This means an arbitrary caller from anywhere will have effective read access to the entire filesystem, regardless of the users' identity.


# Usage

The API has the format

```
http://<hostname>:<port>/getlog
```
With the following query encoded arguments:
- filename  The filename is path relative to `c:\programdata`.  The filename argument is a required argument.  
- count  The number of matching lines to return.  This is an integer value.  A negative or zero value, and all of the lines will be returned.
- match  A case-sensitve match string; if the string appears in a line, then the line is returned.  If this argument is not present or empty, then all lines are returned.
- matchregex  A regular expression match string; use a regular expression to make more complex matches.

Note that match and matchregex are mutually exclusive; providing both returns an error.

An example url would be 
```
"http://localhost:8080/getlog?filename=romeo.txt&count=4
```

which searches the file `c:\programdata\romeo.txt` and returns the most recent (last) 4 lines.

## Example usage

The installer provides a test text file for demonstration purposes.  Once installed, the service can be demonstrated using Powershell's `invoke-webrequest` method to query the service.

### Example 1
```
(iwr -UseBasicParsing -DisableKeepAlive "http://localhost:8080/getlog?filename=romeo.txt&count=4").content
```
will return 

```
[All exit.]
Than this of Juliet and her Romeo.
For never was a story of more woe
Some shall be pardoned, and some punished.
````
(the last 4 lines of the document, in reverse order)

### Example 2

```
(iwr -UseBasicParsing -DisableKeepAlive "http://localhost:8080/getlog?filename=romeo.txt&count=8&match=Romeo").content
```
will return
```
Than this of Juliet and her Romeo.
As rich shall Romeo's by his lady's lie,
[He takes Romeo's letter.]
Where's Romeo's man? What can he say to this?
The noble Paris and true Romeo dead.
Till I conveniently could send to Romeo.
The form of death. Meantime I writ to Romeo
And she, there dead, that Romeo's faithful wife.
```
(the last 8 references to Romeo)

### Example 3 (Search by regex)
```
(iwr -UseBasicParsing -DisableKeepAlive "http://localhost:8080/getlog?filename=romeo.txt&count=8&matchregex=(?i)wherefore").content
```
returns all of the uses of "wherefore", whether capitalized or not:

```
All this is comfort. Wherefore weep I then?
But wherefore, villain, didst thou kill my cousin?
How camest thou hither, tell me, and wherefore?
O Romeo, Romeo, wherefore art thou Romeo?
Why, how now, kinsman? Wherefore storm you so?
```