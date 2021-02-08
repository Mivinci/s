!function () {
  'use-strict'
  const textarea = document.getElementById("textarea");
  const wordCount = document.getElementById("word-count");

  textarea.addEventListener("input", function () {
    wordCount.innerText = `${this.value.length} 字`;
  })

}();