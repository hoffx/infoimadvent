document.addEventListener("DOMContentLoaded", function() {
  renderMathInElement(document.body, {
    displayMode: true,
    delimiters:[
      {left: "$$", right: "$$", display: true},
      {left: "$", right: "$", display: false},
      {left: "\\[", right: "\\]", display: true},
      {left: "\\(", right: "\\)", display: false}
      ]
  });
});