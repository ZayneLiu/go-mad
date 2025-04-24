# Go-MAD: Go Multiple Advanced Data structures and algorithms

Go-MAD is a tool designed to help educate and learn about multiple advanced data structures and algorithms. The tool visualizes the data structures and algorithms in action, allowing users to see how they work and understand their complexities. It also provides interactive features, such as allowing users to input their own data and see how the algorithms perform on it.

## Features

- Visualization of data structures and algorithms
- Interactive TUI (Text User Interface)
- Step-by-step execution of algorithms
- User input for custom data
- Pause, resume, and step through the execution of algorithms

## Supported Data Structures and Algorithms

- Huffman Coding

## Project Structure

- `main.go`: The entry point of the application. This file initializes the TUI and handles user input.
- `tui/`: A directory containing the TUI implementation. This includes custom functions to render data structures and algorithms in the terminal.
- `data_structures/`: A directory containing implementations of the data structures (e.g., Huffman Coding).
- `algorithms/`: A directory containing implementations of the algorithms to be visualized.
- `utils/`: A directory for utility functions that can be used across the project.
- `README.md`: A file providing an overview of the project, instructions for setting up and running the tool, and any other relevant information.
- `.gitignore`: A file specifying which files and directories should be ignored by Git, such as binaries, test outputs, and dependency directories.

## Getting Started

### Prerequisites

- Go 1.24.1 or later

### Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/ZayneLiu/go-mad.git
   cd go-mad
   ```

2. Install dependencies:
   ```sh
   go mod tidy
   ```

### Running the Tool

To run the tool, use the following command:
```sh
go run main.go
```

### Using the Tool

1. Start the tool by running the command above.
2. Follow the on-screen instructions to input your data and visualize the Huffman Coding algorithm.
3. Use the following keybindings to interact with the TUI:
   - `q`: Quit the tool
   - `p`: Pause the visualization
   - `r`: Resume the visualization
   - `s`: Step through the visualization

## Examples

### Visualizing Huffman Coding

1. Input a string to be encoded using Huffman Coding.
2. The tool will display the frequency table of characters and their frequencies.
3. The binary heap used to build the Huffman tree will be visualized.
4. The Huffman tree will be displayed as it is being built.
5. The encoding process will be shown by displaying the binary codes assigned to each character.
6. The decoding process will be visualized by traversing the Huffman tree and reconstructing the original string from the binary input.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue if you have any suggestions or improvements.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [tui-go](https://github.com/marcusolsson/tui-go) for the TUI library
