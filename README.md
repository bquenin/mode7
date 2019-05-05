# Super Nintendo Mode 7 demo

A simple demo reproducing SNES peudo 3D planes.

## What is mode 7?

Mode 7 is a Super Nintendo graphics mode allowing applying an affine transformation (translation, scaling, rotation, etc.) to the background. I strongly recommend [this video explaining mode 7 in depth](https://www.youtube.com/watch?v=3FVN_Ze7bzw).
This mode became famous for its usage in games such as F-Zero and especially Super Mario Kart, allowing to render pseudo 3D planes (technically, pseudo 3D planes on SNES is [a combination of mode 7 and HDMA](https://www.youtube.com/watch?v=K7gWmdgXPgk&feature=youtu.be&t=857))

## How to use

Make sure you have go 1.11 or newer installed.

```shell
go run mode7.go
```
