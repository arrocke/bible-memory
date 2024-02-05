const LETTERS_REGEX = /^[A-Za-z]$/;
const DIACRITIC_REGEX = /[\u0300-\u036f]/g;

/**
 * This constructs a component for a word in this form:
 * <span>
 *   <span>1</span>
 *   <span> "</span>
 *   <span>
 *     <span>w</span>ord
 *   </span>
 *   <span>" </span>
 * </span>
 */
function buildWord({ mode, index, word: data }) {
  const root = document.createElement("span");
  root.className = mode === "learn" ? "text-slate-400" : "text-transparent";

  let number;
  if (data.number) {
    number = document.createElement("span");
    number.innerText = data.number;
    if (index === 0) {
      number.className = "text-black";
    }
    root.appendChild(number);
  }

  let prefix;
  if (data.prefix) {
    prefix = document.createElement("span");
    prefix.innerText = data.prefix;
    root.appendChild(prefix);
  }

  const word = document.createElement("span");

  const firstLetter = document.createElement("span");
  firstLetter.innerText = data.word.charAt(0);
  if (mode === "recall") {
    firstLetter.className = "text-slate-400";
  }
  word.appendChild(firstLetter);

  word.append(data.word.slice(1));
  root.appendChild(word);

  let suffix;
  if (data.suffix) {
    suffix = document.createElement("span");
    suffix.innerText = data.suffix;
    root.appendChild(suffix);
  }

  return {
    root,
    update({ currentIndex, hasHelp, attempts, isCorrect }) {
      const isCurrent = currentIndex === index;

      if (typeof isCorrect === "boolean") {
        if (number) {
          number.className = "";
        }
        if (prefix) {
          prefix.className = "";
        }
        if (suffix) {
          suffix.className = "";
        }
        word.className = "";
        firstLetter.className = "";
        root.className = "";
        if (isCorrect === true && (attempts > 1 || hasHelp)) {
          word.className = "border-b-4 border-orange-400";
        } else if (isCorrect === false) {
          word.className = "border-b-4 border-red-500";
        }
      } else if (hasHelp) {
        word.className = "border-b-4 border-red-500";
        firstLetter.className = "text-slate-400";
      } else if (attempts > 0) {
        word.className = "border-b-4 border-orange-400";
      } else if (isCurrent) {
        if (number) {
          number.className = "text-black";
        }
      }
    },
  };
}

function buildProgress({ size }) {
  const root = document.createElement("div");
  root.className = "flex items-center mb-4";

  const accuracy = document.createElement("div");
  accuracy.className = "mr-2";
  accuracy.innerText = "100%";
  root.appendChild(accuracy);

  const bar = document.createElement("div");
  bar.className =
    "h-2 flex bg-slate-300 items-stretch rounded-full overflow-hidden flex-grow";
  root.appendChild(bar);

  const correctDiv = document.createElement("div");
  correctDiv.className = "bg-green-600";
  correctDiv.style.width = "0%";
  bar.appendChild(correctDiv);

  const incorrectDiv = document.createElement("div");
  incorrectDiv.className = "bg-red-600";
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

function buildRater({ accuracy, onSelect }) {
  const root = document.createElement("form");
  root.className = "mb-8";

  const fieldset = document.createElement("fieldset");
  fieldset.className = "mb-2";
  root.appendChild(fieldset);

  const hardLabel = document.createElement("label");
  const hardInput = document.createElement("input");
  hardLabel.className = "block";
  hardInput.className = "mr-2";
  hardInput.type = "radio";
  hardInput.name = "memory-rating";
  hardInput.value = "2";
  hardInput.required = true;
  hardLabel.appendChild(hardInput);
  hardLabel.appendChild(
    document.createTextNode(
      "Hard - You successfully recalled most of the passage, but you want to review it more frequently"
    )
  );
  fieldset.appendChild(hardLabel);

  const goodLabel = document.createElement("label");
  const goodInput = document.createElement("input");
  goodLabel.className = "flex";
  goodInput.className = "mr-2";
  goodInput.type = "radio";
  goodInput.name = "memory-rating";
  goodInput.value = "3";
  goodInput.required = true;
  goodLabel.appendChild(goodInput);
  goodLabel.appendChild(
    document.createTextNode(
      "Good - You successfully recalled the passage with little help"
    )
  );
  fieldset.appendChild(goodLabel);

  if (accuracy === 1) {
    const easyLabel = document.createElement("label");
    const easyInput = document.createElement("input");
    easyLabel.className = "block";
    easyInput.className = "mr-2";
    easyInput.type = "radio";
    easyInput.name = "memory-rating";
    easyInput.value = "4";
    easyInput.required = true;
    easyLabel.appendChild(easyInput);
    easyLabel.appendChild(
      document.createTextNode(
        "Easy - You successfully recalled the passage with ease and are reviewing it too often"
      )
    );
    fieldset.appendChild(easyLabel);
  }

  const submit = document.createElement("button");
  submit.innerText = "Submit Review";
  submit.className = "button";
  root.appendChild(submit);

  root.addEventListener("submit", (e) => {
    e.preventDefault();
    onSelect?.(parseInt(root["memory-rating"].value));
    root.remove();
  });

  return {
    root,
  };
}

window.Typer = function ({ el: root, words, mode, onComplete }) {
  root.className += " flex flex-col relative";

  const typer = document.createElement("div");
  typer.className =
    "border border-slate-400 rounded p-2 min-h-72 focus-within:outline-2 focus-within:outline outline-blue-600 flex-1 relative";

  const pre = document.createElement("pre");
  pre.className =
    "font-sans whitespace-pre-wrap relative overflow-y-auto absolute w-full h-full";

  const input = document.createElement("input");
  input.className = "absolute opacity-0";

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
    component: buildWord({ mode, index, word }),
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

  input.focus();

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
    const char = e.target.value
      .at(-1)
      .toLowerCase()
      .normalize("NFD")
      .replace(DIACRITIC_REGEX, "");
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
      const accuracy = counts.correct / wordState.length;
      if (mode !== "review") {
        onComplete?.({ accuracy, grade: 0 });
      }
      if (accuracy >= 0.9) {
        const rater = buildRater({
          accuracy,
          onSelect(grade) {
            onComplete?.({ accuracy, grade });
          },
        });
        root.insertBefore(rater.root, root.firstChild);
      } else {
        onComplete?.({ accuracy, grade: 1 });
      }
    } else {
      currentWord = wordState[currentIndex];
      currentWord.component.update({ currentIndex, ...currentWord });
      const newY = Math.max(
        0,
        currentWord.component.root.offsetTop - pre.offsetHeight / 2
      );
      pre.scrollTo(0, newY);
    }

    progress.update(counts);
  });
};
