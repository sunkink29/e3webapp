var scope = 'https://www.googleapis.com/auth/drive.readonly';

var pickerApiLoaded = false;
var oauthToken;

var spreadSheetID = "empty";

var options = null;

// Use the API Loader script to load google.picker and gapi.auth.
function onApiLoad() {
  gapi.load('auth2');
  gapi.load('picker', onPickerApiLoad);
}

function importOnClick() {
  gapi.auth2.authorize({
    client_id: clientId,
    scope: scope
  }, handleAuthResult);
}

function onPickerApiLoad() {
  pickerApiLoaded = true;
  createPicker();
}

function handleAuthResult(authResult) {
  if (authResult && !authResult.error) {
    oauthToken = authResult.access_token;
    createPicker();
  }
}

// Create and render a Picker object for picking user Photos.
function createPicker() {
  if (pickerApiLoaded && oauthToken) {
    var picker = new google.picker.PickerBuilder().
        addView(google.picker.ViewId.SPREADSHEETS).
        setOAuthToken(oauthToken).
        setDeveloperKey(developerKey).
        setCallback(pickerCallback).
        //setRelayUrl("https://e3selectionapp.appspot.com/js/admin/rpc_relay.html").
        build();
    picker.setVisible(true);
  }
}

// A simple callback implementation.
function pickerCallback(data) {
  if (data[google.picker.Response.ACTION] == google.picker.Action.PICKED) {
    var doc = data[google.picker.Response.DOCUMENTS][0];
    spreadSheetID = doc[google.picker.Document.ID];
    showImportOptions();
  }
}

function redirect(data) {
	if (data.url) {
		var strWindowFeatures = "location=yes,height=570,width=520,scrollbars=yes,status=yes";
		var win = window.open(data.url, "_blank", strWindowFeatures);
	}
}

function finishImport() {
  var output = {id: spreadSheetID, options: options};
	postMethod("/admin/importusers", output);
}