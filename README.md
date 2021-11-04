This is a small http server that will read images from a provided path, convert
the image to gray, and then serve back the converted image

Usage
-----

You can start the server like this:

```sh
go run main.go -server 127.0.0.1 -port 8080 -host https://maps.wikimedia.org
```

If you navigate now to `127.0.0.1:8080/osm-intl/1/0/0.png`, the image from
wikimedia will be fetched and the re-served in greyscale
