# CSS Color Variable Creator

A command-line tool that helps you modernize your CSS/SCSS files by automatically converting color values to CSS custom properties (variables).

## Features

- Scans CSS/SCSS files for color values (hex, rgb, rgba)
- Creates CSS custom properties for all found colors
- Generates a new file with all color variables
- Creates a modified version of your input file using the new variables
- Supports both CSS and SCSS files
- Optional conversion of all colors to a specific format (hex, rgb, or rgba)

## Installation

```bash
# Using go install
go install github.com/levinspiekermann/css-color-variable-creator@latest
```

Or download the latest binary from the [releases page](https://github.com/levinspiekermann/css-color-variable-creator/releases).

## Usage

```bash
# Basic usage
css-color-variable-creator create path/to/your/style.css

# Specify output directory
css-color-variable-creator create -d output/dir path/to/your/style.css

# Convert all colors to a specific format (hex, rgb, or rgba)
css-color-variable-creator create -f rgba path/to/your/style.css
```

### Output

The tool generates two files:

1. `{filename}-variables.css`: Contains all color variables
2. `{filename}-with-variables.{ext}`: Your original file modified to use the new variables

### Color Format Conversion

You can use the `-f` or `--format` flag to convert all colors to a specific format:

- `hex`: Convert to hexadecimal format (#RRGGBB or #RRGGBBAA)
- `rgb`: Convert to RGB format (rgb(r, g, b))
- `rgba`: Convert to RGBA format (rgba(r, g, b, a))

For example, to convert all colors to RGBA format:

```bash
css-color-variable-creator create -f rgba style.css
```

## Building from Source

```bash
# Clone the repository
git clone https://github.com/levinspiekermann/css-color-variable-creator.git

# Navigate to the project directory
cd css-color-variable-creator

# Build
go build
```

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
