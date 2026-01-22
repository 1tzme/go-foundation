# go-foundation

This repository is a **foundation monorepo** that aggregates several independent projects originally developed during my studies at the **alem school**.

All projects were migrated from a private school Git platform into a single public GitHub repository **with full preservation of Git history**.

The repository is intended as:
- an archive of completed educational projects
- a single entry point for code review
- a portfolio-style collection of foundational implementations

---

## 📦 Repository Structure

Each project lives in its own directory and remains logically independent.

```
go-foundation/
├── a-library-for-others/
├── bitmap/
├── creditcard/
├── hot-coffee/
├── markov-chain/
├── own-redis/
└── triple-s/
```

---

## 📚 Projects Overview

### `creditcard`
CLI tool for validating, generating, analyzing, and issuing credit card numbers using the Luhn algorithm.

### `markov-chain`
Markov Chain text generator that creates random, coherent text based on input patterns.

### `a-library-for-others`
Library for parsing and handling CSV files with error handling and flexible field extraction.

### `bitmap`
A tool to read, manipulate, and apply filters or transformations to 24-bit BMP image files.

### `triple-s`
Simplified S3-like object storage service with a REST API for managing buckets and objects.

### `hot-coffee`
A RESTful API-based coffee shop management system for handling orders, menu items, and inventory with JSON storage.

### `own-redis`
A custom Redis-like key-value store using UDP protocol for efficient data storage and retrieval.


---

## 🧠 Background

These projects were originally developed as part of the **alem school curriculum**, which emphasizes:
- self-directed learning
- peer review
- problem solving through practice
- writing code close to real systems

The original reference repository for the curriculum structure can be found here:  
- [alem school](https://github.com/alem-platform/foundation)

---

## 📌 Notes

- Each project may use different patterns, tools, or conventions.
- Documentation is located inside each project’s directory.
- Build/run instructions are project-specific.