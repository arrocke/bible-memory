import { KeyboardEvent, useState, Fragment, useEffect } from "react";
import styles from "./VerseTyper.module.css";

const WORD_REGEX = /(\w+(?:'\w+)?)(\W+)(?=\w|$)/g;
const LETTERS_REGEX = /^[A-Za-z]$/;

export interface VerseTyperProps {
  className?: string;
  text: string;
}

interface WordState {
  word: string;
  gap: string;
  attempts: number;
  isCorrect?: boolean;
  hasHelp?: boolean;
}

type WordAction = "correct" | "fail" | "continue" | "help";

export default function VerseTyper({ text, className }: VerseTyperProps) {
  const [words, setWords] = useState<WordState[]>([]);
  useEffect(() => {
    setWords(
      Array.from(text.matchAll(WORD_REGEX), (match) => ({
        word: match[1],
        gap: match[2],
        attempts: 0,
      }))
    );
  }, [text]);

  const currentIndex = words.filter(
    (state) => typeof state.isCorrect === "boolean"
  ).length;
  const currentProgress = words[currentIndex]!;
  const isDone = currentIndex === words.length

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

  function onKeyDown(e: KeyboardEvent) {
    e.preventDefault();
    e.stopPropagation();
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
      default: {
        if (LETTERS_REGEX.test(e.key)) {
          const key = e.key.toLowerCase();
          const firstChar = currentProgress.word[0].toLowerCase();
          attempt(key === firstChar ? "correct" : "fail");
        }
        break;
      }
    }
  }

  return (
    <div
      className={`${className} ${styles.input}`}
      tabIndex={0}
      onKeyDown={onKeyDown}
    >
      {words
        .slice(0, currentIndex + 1)
        .map(({ isCorrect, hasHelp, word, gap, attempts }, i) => {
          const isComplete = typeof isCorrect === "boolean";
          const hadHelp = hasHelp || attempts > 1;
          return (
            <Fragment key={i}>
              <span
                className={
                  isCorrect === false
                    ? styles.incorrectWord
                    : hadHelp
                    ? styles.wordHelp
                    : ""
                }
              >
                {!isComplete && hasHelp ? word[0] : isComplete ? word : null}
              </span>
              {isComplete ? gap : null}
            </Fragment>
          );
        })}
    </div>
  );
}
