partygif
========

A simple command line utility to generate color-rotating animated GIF emojis à la [Party Parrot][party-parrot].

Take your otherwise boring animated GIF emoji...

![gopher-dance](/images/gopher-dance.gif)

...and make it the life of the party!

![party-gopher-dance](/images/party-gopher-dance.gif)

Installation
------------

The easiest way to install `partygif` is to download a pre-built binary from the [releases page][releases].

If you have [Go][golang] installed and want to pull the lastest development version, use the following command.

```
go get -u github.com/ooesili/partygif
```

Usage
-----

```
$ partygif -h
Usage of partygif:
  -black
        add color to a black and white image before color shifting
  -cycles int
        number of color cycles during the GIF (default 1)
  -framerate int
        frame rate in 100ths of seconds for static GIFs (default 10)
  -in string
        input GIF file
  -out string
        output GIF file
  -repeats int
        number of times to repeat GIF before color shifting
```

Example:

```
$ partygif -in gopher-dance.gif -out party-gopher-dance.gif -cycles 3
```

How it works
------------

An animated GIF is composed of a series of frames (pixel tables), each with its own color palette (list of RGB color triplets). Each entry in the pixel table of frame is not an RGB value, but is a index into that frame's color palette. Instead of going after every pixel in the image, this program simply manipulates the color palettes directly.

The color manipulation is done by mapping each RGB value from the color palette to [HCL space][hcl-space] using the super neat [go-colorful][go-colorful] library. It then shifts the hue value of the color by a percentage calculated from the frame number and the `-cycles` flag. After the hue-shift happens, the color is mapped back to the RGB color space and saved back into the palette of the frame.

### Non-animated and short GIFs

For single-frame (A.K.A. non-animated) or short gifs that don't have enough frames to render a smooth color animation, the `-repeats` flag can be used to repeat the GIF a few times before applying the color shifting. For single-frame the `-framerate` flag can be used to customize the frame rate of the final GIF.

### Grayscale Images

Since hue-shifting only works on images that already have color, a `-black` flag exists which will colorize a black and white image before the hue shifting takes place. It does this by maximizing up the red value on every color in every palette on every frame. This essentially turns every black pixel into a red one, which lets the hue-shifting process do its job.

Road map
------------

- [x] repeating short GIFs to add smoother color changes
- [x] static (non-animated) input GIFs
- [ ] PNG input
- [ ] JPEG input
- [x] bringing the party to gray scale images
- [ ] concurrent processing
- [x] pre-built binaries with gox

Thanks
------

Thanks to Egon Elbre for making the dancing gopher GIF used in this example, which came from [this amazing repo][gophers].

Thanks to Lucas Beyer for creating the awesome [go-colorful][go-colorful] library so that I didn't figure out all the math to hue-shift RGB values.


[party-parrot]: http://cultofthepartyparrot.com
[go-colorful]: https://github.com/lucasb-eyer/go-colorful
[hcl-space]: https://en.wikipedia.org/wiki/HCL_color_space
[golang]: https://golang.org/
[gophers]: https://github.com/egonelbre/gophers
[releases]: https://github.com/ooesili/partygif/releases
