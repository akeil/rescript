# reScript

*Handwriting recognition for reMarkable notes.*

This is a small tool to extract text from handwritten notes created on the
[reMarkable](https://remarkable.com/) tablet.

It uses the reMarkable cloud API
and the [MyScript](https://myscript.com/) handwriting recognition ReST API.

## Configuration and Setup
You need to enable "cloud sync" for your remarkable tablet
to make your notes available for this tool.

When first run, *reScript* will ask for a "one time code"
which can be obtained at https://my.remarkable.com/:

```
$ rescript
Enter one time code from https://my.remarkable.com/:
_
```

You also need a [MyScript developer](https://developer.myscript.com/)
account, specifically an `application key` and `HMAC key` which needs to be
added to the configuration file at `~/.config/hwr-conf.yaml`:

```yaml
datadir: /home/USERNAME/.local/share/hwr
cachedir: /home/USERNAME/.cache/hwr
appkey: bbd1419d-aa40-4803-9607-5115c3085de9
hmackey: 33b89262-dde1-4f92-a183-034255db6895
```

The `datadir` and `cachedir` both contain sensitivity values, namely the
authentication token for the reMarkable API, all downloaded notes
and cached handwriting recognition results.

## Usage
Only one use case is supported:

```
$ rescript NAME_OF_NOTE -l LANGUAGE -f FORMAT
```

`NAME_OF_NOTE` is the display name of the notebook you want to convert into
text. It is case-insensitive and supports partial matches.
IF multiple notebooks match, all of them will be converted.

The `LANGUAGE` must be one of the
[languages supported by MyScript](https://developer.myscript.com/docs/interactive-ink/1.4/overview/text-languages/).
The parameter is optional and defaults to `en`.

`FORMAT` specifies the output format. It is either `txt` for plain text
or `md` for markdown.
The parameter is optional and defaults to plain text.

The result is written to a file named after the notebook
in the current directory.

**Example:**

```
$ rescript handwr
… download notebook "Handwriting Recognition"
… recognize handwriting for "Handwriting Recognition"
✓ write "Handwriting Recognition" to "Handwriting Recognition.md"
✓ Done.
```
