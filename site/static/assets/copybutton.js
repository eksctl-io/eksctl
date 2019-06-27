(function() {
    'use strict';

    function flashCopyMessage(el, msg) {
        el.textContent = msg;
        setTimeout(function() {
            el.textContent = "Copy";
        }, 1000);
    }

    function selectText(node) {
        var selection = window.getSelection();
        var range = document.createRange();
        range.selectNodeContents(node);
        selection.removeAllRanges();
        selection.addRange(range);
        return selection;
    }

    function addCopyButton(containerEl) {
        var copyBtn = document.createElement("button");
        copyBtn.className = "highlight-copy-btn";
        copyBtn.textContent = "Copy";

        var codeEl = containerEl.firstElementChild;
        copyBtn.addEventListener('click', function() {
            try {
                var selection = selectText(codeEl);
                document.execCommand('copy');
                selection.removeAllRanges();

                flashCopyMessage(copyBtn, 'Copied!')
            } catch(e) {
                console && console.log(e);
                flashCopyMessage(copyBtn, 'Failed :\'(')
            }
        });

        containerEl.appendChild(copyBtn);
    }

    // Add copy button to code blocks
    var highlightBlocks = document.getElementsByClassName('highlight');
    Array.prototype.forEach.call(highlightBlocks, addCopyButton);

})();
