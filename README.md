# HeLLM: The Language You Never Knew You Didn't Need

Are you tired of your programming language being far too good? Do you find yourself frustrated by fast runtimes, reliable output, and helpful type systems? Introducing **HeLLM** – the only language where *every single statement* is passed through a Large Language Model before it even thinks about running. Why settle for deterministic, logical execution when you could have the full creative chaos of an LLM at your fingertips?

## Features

- **LLM-Powered Everything**  
  Every line, every variable, every print statement – all interpreted by a Large Language Model. Who needs compilers when you have vibes?

- **Revolutionary Function System**  
  Define functions with `fn` and call them with `run`! Because why use normal function calls when you can confuse everyone with custom syntax? Functions even support multiple return values, although only as args as kwargs are far too readable!

- **Pay-as-you-go Interpretation**  
  HeLLM doesn't just interpret your code – it consults an LLM for every single step, ensuring maximum latency and minimum predictability.

- **Nondeterministic Output**  
  Since every statement is filtered through an LLM, you never know what you'll get. Your code might work, or it might write poetry. Isn't that exciting?

- **No Pesky Types**  
  Types are for the weak. In HeLLM, the LLM decides what your variables mean, and sometimes it changes its mind. Function parameters? Just vibes.

- **Blazingly Slow**  
  Why rush? With every statement sent to an LLM, HeLLM ensures you have plenty of time to reflect on your life choices while waiting for your code to finish.

- **Expressive Syntax**  
  Write code that looks almost like English, but not quite enough to be understandable. The LLM will figure it out. Or not.

- **All the Operators You Need!**  
  HeLLM gives you the essentials: `if`, variable assignment, `while` loops, reading input from the CLI, printing to stdout, function definitions with `fn`, function calls with `run`, and even `return` statements! All other operators are useless (trust me, ChatGPT says so, bro).


## Example

```hellm

com "This is a comment, this is ommitted when interpreting";

com "Count to 5, printing each number";
let x = "0";
while "Is x < 5?" {
    let x = "Calculate x + 1";
    print x;
}

com "Query general knowledge";
let city = "The captial city of the UK";
print city;

com "Get the first command line arg";
use name = 0;

com "Define functions and run them from each other";

fn generate_greeting name {
    let greeting = "A greeting for the person called <name>";
    return greeting;
}


fn greet name {
    run msg = generate_greeting name;
    print msg;
}


com "Greet the name, then delete the name from the scope (so no more llm calls see it)";
run greet name;
del name;

com "Use conditional logic, with optional else blocks, and scoping";
let msg = "Empty string";
if "<city> is in europe" {
    let msg = "Exact text: Your city is in europe";
} else {
    let msg = "Exact text: Your city is not in europe";
}
print msg;
```

## Extra Features
- **VSCode Extension Available**  
  Enjoy first-class HeLLM support in Visual Studio Code:
  - Syntax highlighting for all your HeLLM masterpieces.
  - Document formatting to keep your code looking sharp (but beware: extra whitespace is strictly forbidden — this isn't Python, after all).
  - Instant feedback as you type, so you can focus on creative chaos, not code style.

## Installation

Follow these steps to get HeLLM up and running:

1. **Be on Unix**  
   Sorry, Windows users – HeLLM only installs using the makefile on Unix-like systems. You can probably install manually though.

2. **Set Your OPENAI_KEY Environment Variable**  
   HeLLM needs access to an OpenAI API key. Set your `OPENAI_KEY` environment variable in your shell:
   ```
   export OPENAI_KEY=sk-...
   ```
   (Replace `sk-...` with your actual OpenAI API key, so you can track how much money you are ~~wasting~~ enjoying.)

3. **Install Go**  
   You’ll need Go (Golang) installed. If you don’t know how, I’m sure ChatGPT can help.

4. **Install the HeLLM Tool**  
   Open your terminal and run:
   ```
   make tool
   ```
   This will build the HeLLM tool and install it to `/usr/local/bin`.

5. **Install the VSCode/Cursor Extension**  
   Still in your terminal, run:
   ```
   make extension
   ```
   This will copy the extension to your VSCode and Cursor extensions folders.

That's it! You're ready to write the most exciting and expensive code of your life. If you need to see the subcommands of hellm, run `hellm help`.
