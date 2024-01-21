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

function buildProgress({ size }) {
  const root = document.createElement("div");
  root.classList.add("typer-progress");

  const accuracy = document.createElement("div");
  accuracy.classList.add("typer-accuracy");
  accuracy.innerText = "100%";
  root.appendChild(accuracy);

  const bar = document.createElement("div");
  bar.classList.add("typer-progress-bar");
  root.appendChild(bar);

  const correctDiv = document.createElement("div");
  correctDiv.classList.add("typer-progress-correct");
  correctDiv.style.width = "0%";
  bar.appendChild(correctDiv);

  const incorrectDiv = document.createElement("div");
  incorrectDiv.classList.add("typer-progress-incorrect");
  incorrectDiv.style.width = "0%";
  bar.appendChild(incorrectDiv);

  return {
    root,
    update({ correct, incorrect }) {
      correctDiv.style.width = `${(correct / size) * 100}%`;
      incorrectDiv.style.width = `${(incorrect / size) * 100}%`;
      accuracy.innerText = `${((correct / (correct + incorrect)) * 100).toFixed(
        0
      )}%`;
    },
  };
}

window.Typer = function ({ el: root, words, mode, onComplete }) {
  root.classList.add(`typer-${mode}`);

  const typer = document.createElement("div");
  typer.classList.add("typer-input-wrapper");
  const pre = document.createElement("pre");
  pre.classList.add("typer-content");
  const input = document.createElement("input");
  input.classList.add("typer-input");

  typer.addEventListener("click", () => {
    input.focus();
  });

  const wordState = words.map((word, index) => ({
    ...word,
    firstLetter: word.word[0]
      .toLowerCase()
      .normalize("NFD")
      .replaceAll(DIACRITIC_REGEX, ""),
    attempts: 0,
    component: buildWord({ index, word }),
  }));
  const progress = buildProgress({ size: words.length });

  wordState.forEach((word) => {
    pre.appendChild(word.component.root);
  });
  root.appendChild(progress.root);
  typer.appendChild(pre);
  typer.appendChild(input);
  root.appendChild(typer);

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
    const counts = wordState.reduce(
      (counts, word, i) => {
        if (i >= currentIndex) return counts;
        else if (word.isCorrect) counts.correct += 1;
        else counts.incorrect += 1;
        return counts;
      },
      { correct: 0, incorrect: 0 }
    );

    if (currentIndex == wordState.length) {
      input.remove();
      onComplete?.({ accuracy: counts.correct / wordState.length });
    } else {
      currentWord = wordState[currentIndex];
      currentWord.component.update({ currentIndex, ...currentWord });
    }

    progress.update(counts);
  });
};
