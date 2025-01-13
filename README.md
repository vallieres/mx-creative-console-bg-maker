<!-- markdownlint-disable MD002 MD013 MD033 MD041 -->
<h1 align="center">
    <a name="logo" href="https://github.com/vallieres/mx-creative-console-bg-maker">
        <img src="https://github.com/user-attachments/assets/b9a01a70-05a1-4e4c-9f3d-d6c7a99dc504" alt="MX Creative Console Background Maker" width="200">
    </a>
    <br />
    MX Creative Console Background Maker
</h1>
<h4 align="center">
    A command-line tool that splits images into a 3x3 grid for the Logitech MX Creative Console.
</h4>
<div align="center"></div>

<p><font size="3">
        The tool will create 9 PNG files in the same directory as the input image, named with the format `originalname_1.png` through `originalname_9.png`.
    </font></p>

## 🎯 Features

- Resizes images to 378x378px while maintaining aspect ratio
- Intelligently resizes based on the largest dimension
- Crops from the center to preserve image focus
- Splits into 9 equal tiles with proper spacing
- Supports JPEG and PNG input formats
- Outputs individual tiles as PNG files

## ⚡️ Installation

Using Homebrew on MacOS:

```bash
brew tap valliers/mx-creative-console-bg-maker
brew install mx-creative-console-bg-maker
```

Using Go binary:

```bash
go install github.com/vallieres/mx-creative-console-bg-maker/cmd/ccbm@latest
```

## 🛞  Usage

```bash
ccbm <image_path>
```

## License

MIT
