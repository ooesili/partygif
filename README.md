partygif
========

A simple command line utility to generate color-rotating animated gif emoji's à la [Party Parrot][party-parrot].

Take your otherwise boring animated GIF emoji...

![gopher-dance](/images/gopher-dance.gif)

...and make it the life of the party!

![party-gopher-dance](/images/party-gopher-dance.gif)

Usage
-----

```
$ partygif -h
Usage of partygif:
  -cycles int
        number of color cycles during the gif (default 1)
  -in string
        input GIF file
  -out string
        output GIF file
```

Example:

```
$ partygif -in gopher-dance.gif -out party-gopher-dance.gif -cycles 3
```

How it works
------------

An animated GIF is composed of a series of frames (pixel tables), each with its own color palette (list of RGB color triplets). Each entry in the pixel table of frame is not an RGB value, but is a index into that frame's color palette. Instead of going after every pixel in the image, this program simply manipulates the color palettes directly.

The color manipulation is done by mapping each RGB value from the color palette to [HCL space][hcl-space] using the super neat [go-colorful][go-colorful] library. It then shifts the hue value of the color by a percentage calculated from the frame number and the `-cycles` flag. After the hue-shift happens, the color is mapped back to the RGB color space and saved back into the palette of the frame.

Future plans
------------

- [ ] static (non-animated) input gifs
- [ ] png input
- [ ] jpeg input
- [ ] bringing the party to grayscale images
- [ ] concurrent processing

[party-parrot]: http://cultofthepartyparrot.com
[go-colorful]: https://github.com/lucasb-eyer/go-colorful
[hcl-space]: https://en.wikipedia.org/wiki/HCL_color_space
