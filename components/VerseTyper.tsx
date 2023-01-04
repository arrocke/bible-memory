import { KeyboardEvent, useState, Fragment, useEffect, useRef, FormEvent, ReactNode } from "react";
import Button from "./ui/Button";

const WORD_REGEX = /(\d+ )?(\w+(?:'\w+)?)([^A-Za-z']+)(?=\w|\d|$)/g;
const LETTERS_REGEX = /^[A-Za-z]$/;

export interface ProgressUpdate {
  totalWords: number
  wordsComplete: number
  correctWords: number
  correctWordsWithHelp: number
}

export interface VerseTyperProps {
  className?: string;
  text: string;
  mode?: 'review' | 'recall' | 'learn'
  onProgress(progress: ProgressUpdate): void
}

interface WordState {
  prefix: string;
  word: string;
  gap: string;
  attempts: number;
  isCorrect?: boolean;
  hasHelp?: boolean;
}

type WordAction = "correct" | "fail" | "continue" | "help";

function renderWord({ word, isCorrect, hasHelp, attempts, mode = 'review' }: WordState & Pick<VerseTyperProps, 'mode'>): ReactNode {
  const isComplete = typeof isCorrect === "boolean";
  const hadHelp = hasHelp || attempts > 1;

  if (isComplete) {
    if (isCorrect) {
      if (hadHelp) {
        return <span className="bg-yellow-500">{word}</span>
      } else {
        return <span className="inline-block leading-none border-b-2 border-green-400">{word}</span>
      }
    } else {
      return <span className="bg-red-500">{word}</span>
    }
  } else {
    if (mode === 'learn') {
      return <span className="text-gray-500">{word}</span>
    }
    else if (mode === 'recall' || hasHelp) {
      return <span className="text-gray-500">{word[0]}<span className="text-transparent">{word.slice(1)}</span></span>
    } else {
      return null
    }
  }
}

export default function VerseTyper({ text, mode = 'review', className = '', onProgress }: VerseTyperProps) {
  const [words, setWords] = useState<WordState[]>([]);
  useEffect(() => {
    setWords(
      Array.from(text.matchAll(WORD_REGEX), (match) => ({
        prefix: match[1],
        word: match[2],
        gap: match[3],
        attempts: 0,
      }))
    );
  }, [text]);

  const currentIndex = words.filter(
    (state) => typeof state.isCorrect === "boolean"
  ).length;
  const currentProgress = words[currentIndex]!;
  const isDone = currentIndex === words.length

  const wrapper = useRef<HTMLPreElement>(null)
  useEffect(() => {
    setTimeout(() => {
      if (input.current && wrapper.current) {
        const newY = Math.max(0, input.current.offsetTop - wrapper.current.offsetHeight / 2)
        wrapper.current?.scrollTo(0, newY)
      }
    })
    onProgress({
      totalWords: words.length,
      wordsComplete: words.filter(word => typeof word.isCorrect === 'boolean').length,
      correctWords: words.filter(word => word.isCorrect && !word.hasHelp).length,
      correctWordsWithHelp: words.filter(word => word.isCorrect === true).length
    })
  }, [words])

  function attempt(action: WordAction) {
    switch (action) {
      case "correct": {
        setWords((p) => [
          ...p.slice(0, currentIndex),
          {
            ...currentProgress,
            attempts: currentProgress.attempts + 1,
            isCorrect: true,
          },
          ...p.slice(currentIndex + 1),
        ]);
        break;
      }
      case "fail": {
        setWords((p) => [
          ...p.slice(0, currentIndex),
          {
            ...currentProgress,
            attempts: currentProgress.attempts + 1,
          },
          ...p.slice(currentIndex + 1),
        ]);
        break;
      }
      case "continue": {
        setWords((p) => [
          ...p.slice(0, currentIndex),
          {
            ...currentProgress,
            attempts: currentProgress.attempts + 1,
            isCorrect: false,
          },
          ...p.slice(currentIndex + 1),
        ]);
        break;
      }
      case "help": {
        if (mode === 'review') {
          setWords((p) => [
            ...p.slice(0, currentIndex),
            {
              ...currentProgress,
              hasHelp: true,
            },
            ...p.slice(currentIndex + 1),
          ]);
        }
        break;
      }
    }
  }

  const input = useRef<HTMLInputElement>(null)

  function processLetter(char: string) {
    if (LETTERS_REGEX.test(char)) {
      const key = char.toLowerCase();
      const firstChar = currentProgress.word[0].toLowerCase();
      attempt(key === firstChar ? "correct" : "fail");
    }
  }

  function onInput(e: FormEvent<HTMLInputElement>) {
    const char = e.currentTarget.value.at(-1)
    if (char) {
      processLetter(char)
    }
  }

  function onKeyPress(e: KeyboardEvent) {
    if (isDone) return
    switch (e.key) {
      case "/":
      case "?": {
        attempt("help");
        break;
      }
      case "ArrowRight": {
        attempt("continue");
        break;
      }
      default:
        return
    }
    e.preventDefault();
    e.stopPropagation();
  }

  const renderedWords = words
    .slice(0, mode === 'review' ? currentIndex + 1 : undefined)
    .map((data, i) => {
      const { isCorrect, gap, prefix } = data
      return (
        <Fragment key={i}>
          {prefix}
          {renderWord({ ...data, mode })}
          {typeof isCorrect === 'boolean' || mode !== 'review' ? gap : null}
        </Fragment>
      );
    })

  if (!isDone) {
    renderedWords.splice(
      currentIndex,
      0,
      <input
        key="input"
        ref={input}
        className="w-0 focus:outline-none absolute opacity-0"
        onInput={onInput}
        onKeyDown={onKeyPress}
        autoCapitalize="none"
        autoComplete="none"
      />
    )
  }

  return (
    <div className={`${className}`}>
      {
        isDone
          ? null
          : <div className="mb-2">
              <Button
                onClick={() => {
                  attempt('help')
                  input.current?.focus()
                }}
              >
                Hint
              </Button>
              <Button
                className="ml-2"
                onClick={() => {
                  attempt('continue')
                  input.current?.focus()
                }}
              >
                Skip
              </Button>
            </div>
      }
      <pre
        ref={wrapper}
        className="relative focus-within:outline outline-yellow-500 focus-within:border-yellow-500 h-80 overflow-x-auto font-sans whitespace-pre-wrap px-2 py-1 rounded border border-gray-400 shadow-inner select-none"
        tabIndex={isDone ? undefined : -1}
        onFocus={() => input.current?.focus()}
      >
        {renderedWords}
      </pre>
    </div>
  );
}
