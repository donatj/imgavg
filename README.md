# Image Average Generator

Adds a directory of images together, pixel by pixel, and divides by the number of images to generate an average.

Licensed under the [MIT license](http://www.opensource.org/licenses/mit-license.php).

## Operation

```bash
$ imgavg {dir} [outputFile]
```

## Download

See [Releases](https://github.com/donatj/imgavg/releases).

## Usage

	imgavg /path/to/images [outputfilename.png]

## Limitations

- All images must be exactly the same size
- 7.2340173e16 Images Maximum (You'll probably run out of memory first)
