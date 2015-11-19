# GoICNS

GoICNS is a small go utility for converting any PNG image to a full OS X icon
set. It's actually bloody tiny, not small.

When given a PNG image we automatically resize it into each dimension needed by
OS X and create a new ICNS file ready to be saved.

### Why?

I'm fucking bored?  No, actually, I wondered what the format was and whether we
could make it ourselves without calling that shitty `iconutil` command line
tool. Hacky as fuck.

The wikipedia page lists the file format and it was a matter of whipping
somthing up.

### How?

The resizing and file formatting is implemented in `icns.go`; this is a general
purpose library intended to be used within other applications if necessary.

### Show me an example?

An example binary application which accepts a PNG file as input and creates an
ICNS set is included in `main/`.

### One more thing

No.

### License

MIT. Go nuts.
