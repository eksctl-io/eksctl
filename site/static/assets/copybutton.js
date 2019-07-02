(function() {
    'use strict';

    function selectText(node) {
        var selection = window.getSelection();
        var range = document.createRange();
        range.selectNodeContents(node);
        selection.removeAllRanges();
        selection.addRange(range);
        return selection;
    }

    function appendParent(element) {
        var parent = element.parentNode;
        var wrapper = document.createElement('div');
        wrapper.className = "highlight";

        parent.replaceChild(wrapper, element);
        wrapper.appendChild(element);
        return wrapper;
    }

    function addCopyButton(containerEl) {
        var copyBtn = document.createElement("button");
        copyBtn.className = "highlight-copy-btn";
        copyBtn.title = "Copy";


        var preEl = containerEl.parentElement;
        copyBtn.addEventListener('click', function() {
            try {
                var selection = selectText(preEl);
                document.execCommand('copy');
                selection.removeAllRanges();
            } catch(e) {
                console && console.log(e);
            }
        });

        var divEl = appendParent(preEl);
        divEl.appendChild(copyBtn);
    }

    // Add copy button to code blocks
    var highlightBlocks = document.querySelectorAll('pre > code');
    Array.prototype.forEach.call(highlightBlocks, addCopyButton);

})();
