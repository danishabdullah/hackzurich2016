"use strict";

function initGame(questions) {
    var questionField = $("#question"),
        gameProgress = $("#gameProgress"),
        unitField = $("#unit"),
        nextButton = $("#nextQuestion"),
        fieldGroup = $("#boundGroup"),
        minField = $("#lowerBound"),
        maxField = $("#upperBound"),
        idx = 0,
        answers = [];

    function updateView(idx) {
        questionField.html(questions[idx].text);
        gameProgress.html(idx + " / " + questions.length);
        gameProgress.css("width", (idx / questions.length * 100) + "%" );
        unitField.html(questions[idx].unit);
        minField.val("");
        maxField.val("");
    }
    updateView(idx);

    nextButton.click(function() {
        var valid = checkInput(minField, maxField);

        if (!valid) {
            fieldGroup.addClass("has-error");
            return;
        }

        fieldGroup.removeClass("has-error");

        answers[idx] = {
            "question": questions[idx],
            "lower": parseInt(minField.val()),
            "upper": parseInt(maxField.val()),
        }

        if (idx < questions.length - 1) {
            idx++;

            updateView(idx);
            if (idx == (questions.length - 1)) {
                nextButton.html("Finish round");
            } else {
                nextButton.html("Next question");
            }
        } else {
            nextButton.attr("disabled", "disabled");

            $("#formData").val(JSON.stringify(answers));
            $("#gameForm").submit();
        }
    });
}

function isNumber(n) {
  return !isNaN(n) && isFinite(n);
}

function checkInput(minField, maxField) {
    var lower = parseInt(minField.val()),
        upper = parseInt(maxField.val()),
        lowerNumber = isNumber(lower),
        upperNumber = isNumber(upper),
        lowerLessEqual = lower <= upper;

    return lowerNumber && upperNumber && lowerLessEqual;
}
