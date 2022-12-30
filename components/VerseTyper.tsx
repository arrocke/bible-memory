import { KeyboardEvent, useState, Fragment, useEffect, useRef, FormEvent } from "react";
import styles from "./VerseTyper.module.css";

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

export default function VerseTyper({ text, className, onProgress }: VerseTyperProps) {
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

  useEffect(() => {
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
        setWords((p) => [
          ...p.slice(0, currentIndex),
          {
            ...currentProgress,
            hasHelp: true,
          },
          ...p.slice(currentIndex + 1),
        ]);
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
    const value = e.currentTarget.value
    if (value.length > 0) {
      processLetter(value[0])
      e.currentTarget.value = ''
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

  return (
    <div className={`${className} ${styles.wrapper}`}>
      <pre
        tabIndex={isDone ? undefined : -1}
        onFocus={() => input.current?.focus()}
      >
        {words
          .slice(0, currentIndex + 1)
          .map(({ isCorrect, hasHelp, word, gap, prefix, attempts }, i) => {
            const isComplete = typeof isCorrect === "boolean";
            const hadHelp = hasHelp || attempts > 1;
            return (
              <Fragment key={i}>
                {prefix}
                {isCorrect === false ? <span className={styles.incorrectWord}>{word}</span> : null}
                {isCorrect === true && hadHelp ? <span className={styles.wordHelp}>{word}</span> : null}
                {isCorrect === true && !hadHelp ? word : null}
                {!isComplete && hasHelp ? <span className={styles.wordHelp}>{word[0]}</span> : null}
                {isComplete ? gap : null}
              </Fragment>
            );
          })}
          { isDone
            ? null
            : <input
                ref={input}
                className={styles.input}
                onInput={onInput}
                onKeyDown={onKeyPress}
              />
          }
      </pre>
      {
        isDone
          ? null
          : <div>
              <button
                type="button"
                onClick={() => {
                  attempt('help')
                  input.current?.focus()
                }}
              >
                Hint
              </button>
              <button
                type="button"
                onClick={() => {
                  attempt('continue')
                  input.current?.focus()
                }}
              >
                Skip
              </button>
            </div>
      }
    </div>
  );
}
