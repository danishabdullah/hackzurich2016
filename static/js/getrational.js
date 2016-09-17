"use strict";

function initGame(gameID, questions) {
    var questionField = $("#question"),
        unitField = $("#unit"),
        nextButton = $("#nextQuestion"),
        fieldGroup = $("#boundGroup"),
        minField = $("#lowerBound"),
        maxField = $("#upperBound"),
        idx = 0,
        answers = [];

    function updateView(idx) {
        questionField.html(questions[idx].text);
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
            "lower": minField.val(),
            "upper": maxField.val(),
        };

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

            $.ajax({
                type: "POST",
                url: "/game/" + gameID,
                data: answers,
                success: function() {},
                dataType: "json"
            }).fail(function(error) {
                console.log(error);
            }).always(function() {
                nextButton.removeAttr("disabled");
            })
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
