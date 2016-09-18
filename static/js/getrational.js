"use strict";

function initGame(gameID, questions) {
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
        minField.focus();
    }
    updateView(idx);

    function checkReturn(event) {
        var keycode = (event.keyCode ? event.keyCode : event.which);
        if(keycode == '13'){
            nextButtonClick();
        }
    }

    function nextButtonClick() {
        var valid = checkInput(minField, maxField);

        if (!valid) {
            fieldGroup.addClass("has-error");
            return;
        }

        fieldGroup.removeClass("has-error");

        answers[idx] = {
            "question": questions[idx],
            "lower": parseFloat(minField.val()),
            "upper": parseFloat(maxField.val()),
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

            var user = firebase.auth().currentUser;

            $("#formData").val(JSON.stringify({
                "id": gameID,
                "answers": answers,
                "uid": user.uid
            }));
            $("#gameForm").submit();
        }
    }

    minField.keypress(checkReturn);
    maxField.keypress(checkReturn);
    nextButton.click(nextButtonClick);
}

function isNumber(n) {
  return !isNaN(n) && isFinite(n);
}

function checkInput(minField, maxField) {
    var lower = parseFloat(minField.val()),
        upper = parseFloat(maxField.val()),
        lowerNumber = isNumber(lower),
        upperNumber = isNumber(upper),
        lowerLessEqual = lower <= upper;

    return lowerNumber && upperNumber && lowerLessEqual;
}


// firebase stuff

function signIn() {
    console.log("sign in");

    firebase.auth().signInAnonymously().catch(function(error) {
      // Handle Errors here.
      var errorCode = error.code;
      var errorMessage = error.message;
      if (errorCode === 'auth/operation-not-allowed') {
        alert('You must enable Anonymous auth in the Firebase Console.');
      } else {
        console.error(error);
      }
    });
}

function writeUserData(userId) {
  firebase.database().ref('users/' + userId).set({
    exists:true
  });
}

function writeGame(userId, precision) {
    var gameData = {
      uid: uid,
      precision: precision
    };

  var newGameKey = firebase.database().ref().child('games').push().key;

  // Write the new post's data simultaneously in the posts list and the user's post list.
  var updates = {};
  updates['/posts/' + newPostKey] = postData;
  updates['/user-posts/' + uid + '/' + newPostKey] = postData;

  firebase.database().ref('games/' + newGameKey).set(gameData);
}

$.getScript( "https://www.gstatic.com/firebasejs/3.4.0/firebase.js", function() {
  // Initialize Firebase
  var config = {
    apiKey: "AIzaSyDXcUg2dD4jQLfwngfwXGthO3dxPY-jQ_k",
    authDomain: "get-rational.firebaseapp.com",
    databaseURL: "https://get-rational.firebaseio.com",
    storageBucket: "get-rational.appspot.com",
    messagingSenderId: "641794277465"
  };
  firebase.initializeApp(config);

  firebase.auth().onAuthStateChanged(function(user) {
    if (user) {
      // User is signed in.
      var isAnonymous = user.isAnonymous;
      var uid = user.uid;
      console.log("user " + uid + " is signed in.");

      var userId = firebase.auth().currentUser.uid;
      writeUserData(userId);
      firebase.database().ref('/users/' + userId).once('value').then(function(snapshot) {
        if (!snapshot.val()) {
          console.log("user does not exist");
          writeUserData(userId);
        } else {
          console.log("user exists");
        }
      });

    } else {
      // User is signed out.
      console.log("user is signed out.");
      signIn();
    }
});
});
