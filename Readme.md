# BlindX

**BlindX** is a command-line tool for automated Blind XSS testing. Written in Go, it takes a raw HTTP POST request (e.g., copied from Burp Suite), injects your payload into specified parameters or headers, applies various encodings, and reports HTTP responses for each test.

---

## 🚧 Warning

Use **only** on targets for which you have **explicit permission** to test. Unauthorized testing may be illegal and unethical.

---

## ⚙️ Features

* **Raw Request Parsing**: Accepts full HTTP POST requests (headers + body).
* **Parameterized Injection**: Inject into one or more body parameters.
* **Header Injection**: Add or replace headers with payload placeholders.
* **Multiple Encodings**: Choose from 15 single/double/triple encodings (HTML, URL, JS, Unicode, Base64), all together, or none.
* **Batch Execution**: Sends all variants and displays URL + HTTP status code.

---

## 📥 Installation

1. **Prerequisites**

   * Go 1.18 or later

2. **Clone repository**

   ```bash
   git clone https://github.com/yourusername/blindx.git
   cd blindx
   ```

3. **Initialize modules**

   ```bash
   go mod init blindx
   go mod tidy
   ```

4. **Build**

   ```bash
   go build -o blindx main.go
   ```

---

## \$1

## 🛠️ Detailed Workflow

```mermaid
flowchart TD
    A[Start]
    A --> B[Prompt: Paste raw HTTP POST request]
    B --> C[Parse raw text into HTTP Request object]
    C --> D{Parse Success?}
    D -->|No| E[Error: Show parse error & exit]
    D -->|Yes| F[Extract headers & body parameters]

    F --> G[Prompt: Enter first parameter to inject]
    G --> H{More parameters?}
    H -->|Yes| I[Prompt: Enter next parameter]
    I --> H
    H -->|No| J[Collect all parameters]

    J --> K[Prompt: Enter XSS payload]
    K --> L[Prompt: Select encoding option]
    L --> M{Option selected}
    M -->|1–15| N[Apply selected encoding ×n]
    M -->|16| O[Generate all 15 encoding variants]
    M -->|17| P[Use original payload]

    P --> Q[Payload variant list prepared]
    N --> Q
    O --> Q

    Q --> R[Prompt: Additional header injection?]
    R -->|Yes| S[Prompt: Enter header name & value]
    S --> T{More headers?}
    T -->|Yes| S
    T -->|No| U[Headers list prepared]
    R -->|No| U

    U --> V[Initialize HTTP Client]
    V --> W[For each payload variant]
    W --> X[Clone original HTTP Request]
    X --> Y[Inject payload into each parameter]
    Y --> Z[Replace or add headers (with payload placeholders)]
    Z --> AA[Send HTTP request]
    AA --> AB[Receive HTTP response]
    AB --> AC[Log: Request URL + Status Code]
    AC --> AD{More variants?}
    AD -->|Yes| W
    AD -->|No| AE[All tests complete]
    AE --> AF[Exit]
```

---

## 📋 Encoding Details

| Option | Encoding Type     | Variations | Description                 |
| ------ | ----------------- | ---------- | --------------------------- |
| 1–3    | HTML Escape       | ×1, ×2, ×3 | `&lt;` / `&gt;` / etc.      |
| 4–6    | URL Encode        | ×1, ×2, ×3 | `%3C` / `%3E` / etc.        |
| 7–9    | JavaScript Escape | ×1, ×2, ×3 | `\'` / `\"` / `\\` escapes  |
| 10–12  | Unicode Escape    | ×1, ×2, ×3 | `\u003C` / `\u003E` / etc.  |
| 13–15  | Base64 Encode     | ×1, ×2, ×3 | `PHNjcmlwdD4=` etc.         |
| 16     | All Variants      | 15 total   | Runs all above in sequence  |
| 17     | None              | —          | Original payload unmodified |

---

## 📄 License

This project is released under the **MIT License**.

---

*Developed by Prog & Contributors*
