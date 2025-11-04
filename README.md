# Fail2Ban Log Parser

## About

This is a mini weekend project i did as one of my first steps in learning golang. The binary takes a f2b log file as input and a output directory as output, in which it creates a parsed.json and state.json file. Parsed.json contains the logs in a structured file and state.json is the metadata used to track at which point the parser last checked and if the file rolled over (log files tend to wipe after a certain amount of rows or create copies of themselfes while wiping their contents.)

## Usage

To build the project run

```bash
make build
```

This will create a binary: `bin/parser`

Once you run the binary you will get the usage menu:
```
Usage: ./bin/parser --destDir=<dir> --source=<file>
  -d string
    	Destination Directory (shorthand)
  -destDir string
    	Destination Directory
  -s string
    	Source Log File (shorthand)
  -source string
    	Source Log File
```

Lets say you have the following structure:

```
├── fail2ban/
├── parser
└── raw-logs/
    └── fail2ban.log
```

You can run the parser:
```
./parser -destDir=./fail2ban -source=./raw-logs/fail2ban.log
```

The parser will read the source file and create `parsed.json`and `state.json` in your destination directory. _*If destination directory does not exist, the parser will try to create it._

`parsed.json` will have the following structure:

```
[
   {
      "timestamp": "2025-10-19 00:00:01,810",
      "handler": "fail2ban.server",
      "level": "INFO",
      "source": "",
      "ipAddress": "",
      "message": ""
   },
...
]
```

You can import this in grafana or any other tool to analyse further.

`state.json` will look like this:

```
{
  "offset": 3703746,
  "size": 6079112
}
```

Offset is the size in Bytes of how far the fail2ban.log file was read when the command was run the last time.

Size is the tracker of the parsed log file size in Bytes. - I have implemented this for myself to have an overview, maybe one day I will write rollover logic based on this.


## More

To read more about this project and how i made a Grafana dashboard you can read it on [My Blog](https://agres.online/blog/bruteforce)


## Developed with

- make
- go 1.25.3
- linux
