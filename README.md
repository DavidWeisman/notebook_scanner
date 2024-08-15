# Notebook Scanner

## Overview

This project provides a tool for scanning code repositories to analyze notebook files. The scanner processes notebooks and generates reports on their contents.

## Prerequisites

It is recommended to use a virtual environment for managing dependencies.

## Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/DavidWeisam/notebook_scanner.git
    ```

2. **Navigate to the `scanner` directory:**

    ```bash
    cd notebook_scanner/scanner
    ```

3. **Install the requirements:**

    ```bash
    ./install_all.sh
    ```

## Usage

To run the scan on a repository, use the following command:

```bash
go run main.go <repo-url> <path-to-folder-to-scan>
```

## Results

The scan results will be saved in the following locations:

- `scan_reports` folder
- `scanner` folder

