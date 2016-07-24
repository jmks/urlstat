# URLstat

Prints out the HTTP status of URLs it extracts from a list of files.

It's been an exercise of experimenting with Go routines, but is a useful tool to check links in files.

# Example Usage
```
$ urlstat --list file-of-urls
http://www.example.com
http://localhost:9000/200
http://localhost:9000/404
http://localhost:9000/500

$ urlstat file-of-urls
or
$ cat file-of-urls | urlstat
404 Not Found : http://localhost:9000/404
200 OK : http://localhost:9000/200
500 Internal Server Error : http://localhost:9000/500
200 OK : http://www.example.com

$ urlstat -ok file-of-urls
200 OK : http://localhost:9000/200
200 OK : http://www.example.com

$ urlstat -no-ok file-of-urls
404 Not Found : http://localhost:9000/404
500 Internal Server Error : http://localhost:9000/500
```

Note: -no-ok trumps -ok
```
$ urlstat --ok --no-ok file-of-urls
404 Not Found : http://localhost:9000/404
500 Internal Server Error : http://localhost:9000/500
```

## TODO
- re-add tests!!
- move concurrency to slowest part of the pipeline
- add context to matches (filename and line number)
- group output by filename, order by line number
- match URNs (URI without scheme) if its TLD is [valid](http://data.iana.org/TLD/tlds-alpha-by-domain.txt)
- Able to update TLDs (user-defined whitelist?)
- Skip scanning binary files
