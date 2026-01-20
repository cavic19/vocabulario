# Vocabulario
A vocabulary learning app with spaced repetition. Words you struggle with appear more frequently.

## Setup
1. Create a `data/` directory with CSV files (format: `word1;word2`)
2. Run: `go run . -dir data`
3. Answer the prompts—the app tracks your progress and focuses on difficult words

## Features
- Multiple lesson files (e.g., `leccion_1.csv`, `leccion_2.csv`)
- Spaced repetition based on success/failure tracking
- Statistics saved automatically to `stats.json`
- Accent-insensitive comparison (optional)
- Multiple correct answers per word—use separate rows in your CSV for each pairing (e.g., `A;B` and `A;C` to accept B or C as correct answers for A)

## TODO
- [ ] Integrate with AI agent to generate contextual sentences for translation using the user's most-struggled words


