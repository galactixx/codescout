<p align="center">
  <img src="/docs/logo.png" alt="codescout logo" width="65%"/>
</p>

**CodeScout** is a Go package and CLI tool for analyzing and extracting structured representations of functions, methods, and structs from Go source files.

---

### ✍️ Exported Functions

#### `ScoutFunction(path string, config FuncConfig) (*FuncNode, error)`
Returns the first function matching the configuration in the provided file path.

#### `ScoutFunctions(path string, config FuncConfig) ([]*FuncNode, error)`
Returns all matching functions based on the configuration.

#### `ScoutStruct(path string, config StructConfig) (*StructNode, error)`
Returns the first struct matching the configuration in the given file path.

#### `ScoutStructs(path string, config StructConfig) ([]*StructNode, error)`
Returns all structs matching the provided configuration.

#### `ScoutMethod(path string, config MethodConfig) (*MethodNode, error)`
Returns the first method matching the configuration.

#### `ScoutMethods(path string, config MethodConfig) ([]*MethodNode, error)`
Returns all methods that match the given configuration.

### ⚖️ Configuration Types

#### `FuncConfig`
Defines filters for scouting functions including:
- Name
- Parameter and return types
- Match options: `Exact`, `NoParams`, `NoReturn`

#### `MethodConfig`
Used to find specific methods, with support for:
- Receiver type
- Pointer receiver flag
- Accessed fields and called methods
- Match options: `Exact`, `NoParams`, `NoReturn`, `NoFields`, `NoMethods`

#### `StructConfig`
Defines search criteria for structs:
- Field name and type matches
- `Exact` and `NoFields` options

---

## ⚖️ CLI Usage

CodeScout also offers a full-featured CLI using [Cobra](https://github.com/spf13/cobra).

### 🔢 Function Command
```bash
codescout func [path] [flags]
```
- `--name`, `-n`: Function name
- `--params`, `-p`: Function parameters
- `--return`, `-r`: Return types
- `--no-params`, `-s`: Expect no parameters
- `--no-return`, `-u`: Expect no return values
- `--exact`, `-x`: Match criteria exactly
- `--output`, `-o`: Output format (`definition`, `body`, `signature`, etc.)

### 🎓 Method Command
```bash
codescout method [path] [flags]
```
- `--name`, `-n`: Method name
- `--receiver`, `-m`: Receiver type
- `--pointer`, `-t`: Whether it's a pointer receiver
- `--fields`, `-f`: Fields accessed
- `--methods`, `-c`: Methods called
- `--no-fields`, `-d`: Must not access struct fields
- `--no-methods`, `-e`: Must not call struct methods
- All function flags also apply

### 💼 Struct Command
```bash
codescout struct [path] [flags]
```
- `--name`, `-n`: Struct name
- `--fields`, `-f`: Fields to match
- `--no-fields`, `-s`: Struct must have no fields
- `--exact`, `-x`: Match fields exactly
- `--output`, `-o`: Output format (`definition`, `body`, etc.)

### 💡 Verbose Output
All commands support the `--verbose`, `-v` flag to list **all** matches instead of just the first.

---

## 📅 Example
```bash
codescout func ./example.go -name=SomeFunc -params=input:string -return=error -output=signature
```

---

## 🚀 Getting Started

### 💼 Use as a Package

```bash
go get github.com/galactixx/codescout@latest
```

Then import in your Go code:

```go
import "github.com/galactixx/codescout"
```

---

### 🛠️ Install & Use the CLI

```bash
go install github.com/galactixx/codescout/cmd/codescout@latest
```

Then run:

```bash
codescout
```

---

## 🔮 **Future Features**

- Support for scouting interfaces and their methods.
- Ability to search for structs that implement specific interfaces via MethodConfig matching.
- Color-coded Go syntax highlighting in CLI output.
- Smart path resolution and recursive file scanning.
- Integration with gopls for enhanced analysis.

---

## 🤝 **License**

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

---

## 📞 **Contact**

If you have any questions or need support, feel free to reach out by opening an issue on the [GitHub repository](#).
