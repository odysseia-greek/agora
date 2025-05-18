# Plato - Πλάτων

> χαλεπὰ τὰ καλά - good things are difficult to attain

Common layer for all odysseia-greek applications. Plato serves as the foundational utility package that provides shared functionality across the entire odysseia-greek ecosystem.

## Overview

Plato is designed as a modular utility package that provides various components to standardize and simplify common operations across the odysseia-greek project. It includes utilities for:

- Configuration management
- Certificate generation and TLS management
- Progress tracking for educational applications
- Random number generation
- Logging
- Service middleware
- Data transformation
- And more

## Components

### Certificates
Provides functionality for generating and managing TLS certificates, including:
- Certificate Authority (CA) creation
- Key pair and certificate generation
- PEM encoding of certificates and private keys

### Config
Central configuration hub containing constants, environment variable definitions, and default values used throughout the odysseia-greek ecosystem.

### Generator
Utilities for generating various types of data.

### Helpers
Common helper functions used across applications.

### Logging
Standardized logging functionality.

### Middleware
Service middleware components for HTTP services.

### Models
Shared data models and structures.

### Progress
Functionality for tracking progress in learning sessions, particularly for language learning applications:
- Word progress tracking
- Session progress management
- Accuracy calculation
- Segment completion tracking

### Randomizer
Utilities for generating random numbers with secure seeding:
- Base-zero random number generation (0 to n-1)
- Base-one random number generation (1 to n)

### Service
Common service-related utilities.

### TLS Manager
Tools for managing TLS configurations.

### Transform
Data transformation utilities.

## Installation

```bash
go get github.com/odysseia-greek/agora/plato
```

## Usage

### Randomizer Example

```go
package main

import (
    "fmt"
    "github.com/odysseia-greek/agora/plato/randomizer"
)

func main() {
    // Create a new randomizer client
    r, err := randomizer.NewRandomizerClient()
    if err != nil {
        // Handle error
    }

    // Generate a random number from 0 to 9
    randomNumber := r.RandomNumberBaseZero(10)
    fmt.Printf("Random number (0-9): %d\n", randomNumber)

    // Generate a random number from 1 to 10
    randomNumberOne := r.RandomNumberBaseOne(10)
    fmt.Printf("Random number (1-10): %d\n", randomNumberOne)
}
```

### Progress Tracking Example

```go
package main

import (
    "fmt"
    "github.com/odysseia-greek/agora/plato/progress"
)

func main() {
    // Create a progress tracker
    tracker := &progress.ProgressTracker{
        Sessions: make(map[string]map[string]map[string]*progress.WordProgress),
    }

    // Initialize words for a segment
    sessionId := "user123"
    segmentKey := "lesson1"
    greekWords := []string{"ἄνθρωπος", "θεός", "λόγος"}

    tracker.InitWordsForSegment(sessionId, segmentKey, greekWords)

    // Record a word play
    tracker.RecordWordPlay(sessionId, segmentKey, "ἄνθρωπος", "human")

    // Record an answer result
    tracker.RecordAnswerResult(sessionId, segmentKey, "ἄνθρωπος", true)

    // Get playable words
    unplayed, unmastered := tracker.GetPlayableWords(sessionId, segmentKey, 3)

    fmt.Printf("Unplayed words: %v\n", unplayed)
    fmt.Printf("Unmastered words: %v\n", unmastered)
}
```

### Certificate Generation Example

```go
package main

import (
    "fmt"
    "github.com/odysseia-greek/agora/plato/certificates"
)

func main() {
    // Create a certificate generator
    generator := &certificates.CertificateGenerator{
        Organizations: []string{"odysseia-greek"},
        CaValidity:    365, // validity in days
    }

    // Initialize the CA
    err := generator.InitCa()
    if err != nil {
        // Handle error
    }

    // Generate a key and certificate set
    hosts := []string{"example.com", "www.example.com"}
    validityDays := 365

    certPEM, keyPEM, err := generator.GenerateKeyAndCertSet(hosts, validityDays)
    if err != nil {
        // Handle error
    }

    fmt.Printf("Certificate generated for hosts: %v\n", hosts)
    fmt.Printf("Certificate length: %d bytes\n", len(certPEM))
    fmt.Printf("Private key length: %d bytes\n", len(keyPEM))
}
```
