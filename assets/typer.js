const LETTERS_REGEX = /^[A-Za-z]$/;
const DIACRITIC_REGEX = /[\u0300-\u036f]/g;

/**
 * This constructs a component for a word in this form:
 * <span>
 *   <span class="typer-prefix">prefix</span>
 *   <span class="typer-word">
 *     <span class="typer-first-letter">w</span>ord
 *   </span>
 *   <span class="typer-suffix">suffix</span>
 * </span>
 */
function buildWord({ index, word: data }) {
  const root = document.createElement("span");
  if (index === 0) {
    root.classList.add("typer-current");
  }

  if (data.prefix) {
    const prefix = document.createElement("span");
    prefix.className = "typer-prefix";
    prefix.innerText = data.prefix;
    root.appendChild(prefix);
  }

  const word = document.createElement("span");
  word.className = "typer-word";

  const firstLetter = document.createElement("span");
  firstLetter.className = "typer-first-letter";
  firstLetter.innerText = data.word.charAt(0);
  word.appendChild(firstLetter);

  word.append(data.word.slice(1));
  root.appendChild(word);

  if (data.suffix) {
    const suffix = document.createElement("span");
    suffix.className = "typer-suffix";
    suffix.innerText = data.suffix;
    root.appendChild(suffix);
  }

  return {
    root,
    update({ currentIndex, hasHelp, attempts, isCorrect }) {
      const isCurrent = currentIndex === index;

      if (isCorrect === true) {
        root.className =
          attempts > 1 || hasHelp ? "typer-almost" : "typer-correct";
      } else if (isCorrect === false) {
        root.className = "typer-incorrect";
      } else if (hasHelp) {
        root.className = "typer-hint typer-current";
      } else if (attempts > 0) {
        root.className = "typer-attempt typer-current";
      } else if (isCurrent) {
        root.className = "typer-current";
      } else {
        root.className = "";
      }
    },
  };
}

window.Typer = function ({ el: root, words, mode, onProgress }) {
  root.classList.add(`typer-${mode}`);
  const pre = document.createElement("pre");
  const input = document.createElement("input");

  const wordState = words.map((word, index) => ({
    ...word,
    firstLetter: word.word[0]
      .toLowerCase()
      .normalize("NFD")
      .replaceAll(DIACRITIC_REGEX, ""),
    attempts: 0,
    component: buildWord({ index, word }),
  }));

  wordState.forEach((word) => {
    pre.appendChild(word.component.root);
  });
  root.appendChild(pre);
  root.appendChild(input);

  let currentIndex = 0;
  let currentWord = wordState[currentIndex];

  if (mode === "review") {
    input.addEventListener("keydown", (e) => {
      if (e.key === "Enter") {
        currentWord.hasHelp = true;
        currentWord.component.update({ currentIndex, ...currentWord });
        e.preventDefault();
        e.stopPropagation();
      }
    });
  }

  input.addEventListener("input", (e) => {
    const char = e.target.value.at(-1);
    e.target.value = "";
    if (!char || !LETTERS_REGEX.test(char)) {
      return;
    }

    currentWord.attempts += 1;

    if (char.toLowerCase() !== currentWord.firstLetter) {
      navigator.vibrate(100);
      currentWord.component.update({ currentIndex, ...currentWord });
      return;
    }

    currentWord.isCorrect = !currentWord.hasHelp;
    currentWord.component.update({ currentIndex, ...currentWord });

    currentIndex += 1;
    if (currentIndex == wordState.length) {
      input.remove();
    } else {
      currentWord = wordState[currentIndex];
      currentWord.component.update({ currentIndex, ...currentWord });
    }

    const correctCount = wordState.reduce(
      (count, word, i) =>
        i < currentIndex && word.isCorrect ? count + 1 : count,
      0
    );
    onProgress?.({
      progress: currentIndex / wordState.length,
      accuracy: correctCount / currentIndex,
    });
  });
};
