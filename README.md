<p align="center">
 <img src="assets/banner.jpg" width="400">
</p>

<div align="center">

  <a href="">![Tests](https://github.com/leonidasdeim/log/actions/workflows/go.yml/badge.svg)</a>
  <a href="">![Code Scanning](https://github.com/leonidasdeim/log/actions/workflows/codeql.yml/badge.svg)</a>
  <a href="https://codecov.io/gh/leonidasdeim/log" > 
    <img src="https://codecov.io/gh/leonidasdeim/log/branch/main/graph/badge.svg?token=3275GV3OGX"/> 
  </a>
  <a href="">![Report](https://goreportcard.com/badge/github.com/leonidasdeim/log)</a>
  <a href="">![Release](https://badgen.net/github/release/leonidasdeim/log)</a>
  <a href="">![Releases](https://badgen.net/github/releases/leonidasdeim/log)</a>
  
</div>

# log

A versatile and modular logging library designed specifically for Go applications. With `log`, you can effortlessly manage and organize your logs in a way that suits your modular application structure. Whether you need a global logger for your application or specific local loggers for individual modules, `log` has got you covered.


```bash
go get github.com/leonidasdeim/log
```

# Features

**Modular Logging:** `log` is tailored for modular applications, allowing you to create both global and local loggers. This flexibility empowers you to manage logs efficiently across different parts of your application while maintaining global properties.

**Global and Local Logging:** Enjoy the best of both worlds. Use a main logger for global logging requirements, while simultaneously creating local loggers with distinct properties for individual modules or components.

**Customizable Logging Levels:** `log` supports customizable logging levels, enabling you to fine-tune the verbosity of your logs. Choose from a range of logging levels such as DEBUG, INFO, WARNING, ERROR, and FATAL.

**Formatted Logging:** Format your log messages the way you want. `log` supports flexible log message formatting to suit your needs.

**Concurrency-Safe:** Built to handle concurrent access safely, `log` ensures that your logs won't get tangled when multiple goroutines are writing to the same logger.

**Extensible:** Easily extend `log` with custom log output targets or adjust its behavior to fit your specific application requirements.

# Usage

Please check usage examples in `examples/`