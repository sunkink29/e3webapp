var regularSheetId = '1Q53QVvSGb4wb7KbV1eJRCusPaAmeGU3a2LTvtCDeh2k';
var debugSheetId = '1OPCahVswUJuDtNjX-KrCpQxjaTMca8vUgcuPpHvUJx8';
var currentWeek = false;
var ss;
var curUserRow;
var teacherList;
var isTeacherCache;
var studentList;
var adminList;
var testCol = 'H';
var tStartCol = 'A';
var sStartCol = 'F';
var emailFactor = 1;
var nameFactor = 2;
var block1Factor = 4;
var block2Factor = 5;
var tBlockStart = 3;
var bLength = 6;
var bStatus = 5;
var bCurSize = 3;
var bMaxSize = 4;

function doGet(e) {
  var cache = CacheService.getUserCache();
  cache.put('debug', e.parameter.debug === "true");
  
  var t = HtmlService.createTemplateFromFile('index');
  t.data = {teacher: false, admin: false};
  t.data.teacher = isTeacher(); 
  t.data.admin = isAdmin();
  
  return t.evaluate().setTitle('E3 App')
  .setXFrameOptionsMode(HtmlService.XFrameOptionsMode.ALLOWALL);
}

function addLetter(letter, factor) {
  var charCode = letter.charCodeAt();
  return String.fromCharCode(charCode+factor);
}

function openSpreadSheet() {
  if (ss == null) {
    var cache = CacheService.getUserCache();
    var useDebug = cache.get('debug') === "true";
    var sheetId = useDebug?debugSheetId : regularSheetId;
    ss = SpreadsheetApp.openById(sheetId);
  }
}

function getStudentList() {
  if (studentList == null) {
    openSpreadSheet();
    var sheet = ss.getSheetByName("student data " + (currentWeek===true?"current":"new"));
    studentList = sheet.getRange('A2:'+ addLetter('A', sheet.getLastColumn()) + sheet.getLastRow()).getValues();
   }
}

function getTeacherList() {
  if (teacherList == null) {
      openSpreadSheet();
      var sheet = ss.getSheetByName("teacher data " + (currentWeek===true?"current":"new"));
      teacherList = sheet.getRange('A2:'+ addLetter('A', sheet.getLastColumn()) + sheet.getLastRow()).getValues();
   }
}

function getUserList() {
  openSpreadSheet();
  var sheet = ss.getSheetByName("permission data");
  return sheet.getRange('A2:'+ addLetter('A', sheet.getLastColumn()) + sheet.getLastRow()).getValues();
}

function getPermissions() {
  var permissions = {isTeacher: false, isAdmin: false};
  var cache = CacheService.getUserCache();
  permissions.isTeacher = cache.get('isTeacher');
  permissions.isAdmin = cache.get('isAdmin');
  
  if (permissions.isTeacher == null || permissions.isAdmin == null) {
    openSpreadSheet();
    var sheet = ss.getSheetByName("permission data");
    var permissionList = sheet.getRange('A2:'+ addLetter('A', sheet.getLastColumn()) + sheet.getLastRow()).getValues();
    var currentUser = Session.getActiveUser();
    permissionList.forEach(function(item, index) {
      if (item[1] == currentUser) {
        permissions.isTeacher = item[3];
        permissions.isAdmin = item[4];
      }
    });
    cache.put('isTeacher', permissions.isTeacher);
    cache.put('isAdmin', permissions.isAdmin);
  }
  return permissions;
}

function isTeacher() { 
  if (isTeacherCache == null) {
    isTeacherCache = getPermissions().isTeacher;
  }
  return (isTeacherCache === 'true') || isTeacherCache === true;
}

function isAdmin() {
  isAdmin = getPermissions().isAdmin;
  return (isAdmin === 'true') || isAdmin === true;
}

function setCurUserRow() {
  var cache = CacheService.getUserCache();
  curUserRow = cache.get('curUserRow');
  if (curUserRow == null) {
    if (isTeacher() == true) {
      getTeacherList();
      var currentUser = Session.getActiveUser();
      teacherList.forEach(function(item, index) {
        if (item[1] == currentUser) {
          curUserRow = item[0]+2;
          cache.put('curUserRow', curUserRow);
        }
      })
    } else {
      getStudentList();
      var currentUser = Session.getActiveUser();
      studentList.forEach(function(item, index) {
        if (item[1] == currentUser) {
          curUserRow = item[0]+2;
          cache.put('curUserRow', curUserRow);
        }
      })
    }
  }
}

function getStudents(showAll, grade, currentWeek, teacherIndex) {
  if (isTeacher() || isAdmin()) {
    this.currentWeek = currentWeek
    var block1List = [];
    var block2List = [];
    getStudentList();
    getTeacherList();
    setCurUserRow();
    var curUserRow = this.curUserRow;
    
    if (teacherIndex != undefined && isAdmin()) {
      curUserRow = teacherIndex + 2;
    }
    
    studentList.forEach(function(item, index) {
      if (!(grade == undefined) == (item[3] == grade)) {
        if (showAll == true) {block1List.push(item);}
        else {
          if (item[5] == curUserRow-1) {block2List.push(item);}
          if (item[4] == curUserRow-1) {block1List.push(item);}
        }
      }
    })
    studentList.forEach(function(item, index) {
       item[block1Factor] = teacherList[parseInt(item[block1Factor])-1];
       item[block2Factor] = teacherList[parseInt(item[block2Factor])-1];
    })
    if (showAll != true) {
      return [block1List, block2List];
    } else {
      return block1List;
    }
  }
}

function getTeachers(current, currentWeek) {
  this.currentWeek = currentWeek;
  var output = [];
  getTeacherList();
  if (current == true) {
    setCurUserRow();
    openSpreadSheet();
    var ssheet = ss.getSheetByName("student data " + (currentWeek===true?"current":"new"));
    var currentTeachers = [ssheet.getRange('E'+curUserRow).getValue(),ssheet.getRange('F'+curUserRow).getValue()];
    currentTeachers.forEach(function(item, index) {
      teacher = teacherList[item-1];
      output.push(teacher);
    });
  } else {
    output = teacherList;
  }
  return output
}

function setClass(block, teacherId, currentWeek) {
  this.currentWeek = currentWeek;
  setCurUserRow();
  getStudentList();
  getTeacherList();
  var currentBlockOpen;
  if (studentList[curUserRow-2][block1Factor+block] >= 1) {
    currentBlockOpen = teacherList[studentList[curUserRow-2][block1Factor+block]-1][tBlockStart+bLength*(block)+bStatus] == 'open';
  } else {
    currentBlockOpen = true
  }
  var newBlockOpen = teacherList[teacherId][tBlockStart+bLength*(block)+bStatus] == 'open';
  var newBlockFull = teacherList[teacherId][tBlockStart+bLength*(block)+bCurSize] >= teacherList[teacherId][tBlockStart+bLength*(block)+bMaxSize];
  var message = {succeed: false};
  if (currentBlockOpen && newBlockOpen && !newBlockFull) {
    openSpreadSheet();
    var sheet = ss.getSheetByName("student data " + (currentWeek===true?"current":"new"));
    sheet.getRange(addLetter('A', block+4)+curUserRow).setValue(teacherId+1);
    message.succeed = true;
  } else {
    if (!currentBlockOpen) {
      message.error = "Current class is a closed session";
    } else if (!newBlockOpen) {
      message.error = "New class is a closed session";
    } else if (newBlockFull) {
      message.error = "New class is Full";
    }
  }
  return message;
}

function addStudentToClass(studentId, block, currentWeek) {
  var message = {succeed: false, error: ""};
  if (isTeacher()) {
    this.currentWeek = currentWeek;
    getStudentList();
    getTeacherList();
    setCurUserRow();
    var teacherId = curUserRow - 2;
    var studentClassOpen =  true;
    if (teacherList[studentList[studentId][block1Factor+block]-1] != null) {
      teacherList[studentList[studentId][block1Factor+block]-1][tBlockStart+bLength*(block)+bStatus] == 'open';
    }
    var classFull = teacherList[teacherId][tBlockStart+bLength*(block)+bCurSize] >= teacherList[teacherId][tBlockStart+bLength*(block)+bMaxSize];
    if (studentClassOpen && !classFull) {
      openSpreadSheet();
      var sheet = ss.getSheetByName("student data " + (currentWeek===true?"current":"new"));
      sheet.getRange(addLetter('A', block+4)+(studentId+2)).setValue(curUserRow-1);
      message.succeed = true;
    } else if (!studentClassOpen) {
      message.error = "Student previous class is a closed session";
    } else {
      message.error = "Class full";
    }
  } else {
    message.error = "Teacher permissions required";
  }
  return message;
}

function removeStudentFromClass(studentId, block, currentWeek) {
  var message = {succeed: false, error: ""};
  if (isTeacher()) {
    this.currentWeek = currentWeek;
    openSpreadSheet();
    var sheet = ss.getSheetByName("student data " + (currentWeek===true?"current":"new"));
    sheet.getRange(addLetter('A', block+4)+(studentId+2)).setValue('');
    message.succeed = true;
  } else {
    message.error = "Teacher permissions required";
  }
  return message;
}

function getBlockInfo(currentWeek) {
  var message = {succeed: false, error: "", blockInfo:[]}
  if (isTeacher()) {
    this.currentWeek = currentWeek;
    getTeacherList();
    setCurUserRow();
    message.blockInfo = teacherList[curUserRow - 2].slice(tBlockStart,tBlockStart+bLength*2);
    message.succeed = true;
  } else {
    message.error = "invalid permissions";
  }
  return message;
}

function editBlockInfo(blockInfo, currentWeek) {
  var message = {succeed: false, error: ""}
  if (isTeacher()) {
    setCurUserRow();
    this.currentWeek = currentWeek;
    openSpreadSheet();
    var sheet = ss.getSheetByName("teacher data " + (currentWeek===true?"current":"new"));
    var blockCells = sheet.getRange(curUserRow, tBlockStart+bLength*blockInfo.blockNum+1, 1, bLength+1);
    blockCells.getCell(1, 1).setValue(blockInfo.subject);
    blockCells.getCell(1, 2).setValue(blockInfo.description);
    blockCells.getCell(1, 3).setValue(blockInfo.roomNumber);
    blockCells.getCell(1, 4).setValue(blockInfo.curSize);
    blockCells.getCell(1, 5).setValue(blockInfo.maxSize);
    blockCells.getCell(1, 6).setValue(blockInfo.status);
    message.succeed = true;
  } else {
    message.error = "invalid permissions";
  }
  return message;
}

function addUser(user) {
  if (isAdmin()) {
    openSpreadSheet();
    var sheet = ss.getSheetByName("Permission data");
    var newRow = sheet.getRange(sheet.getLastRow()+1, 1, 1, 6);
    var userRow = "";
    newRow.getCell(1, 1).setFormula("row()-2");
    newRow.getCell(1, 2).setValue(user.email);
    newRow.getCell(1, 3).setValue(user.name);
    newRow.getCell(1, 4).setValue(user.isTeacher.toLowerCase() == 'true');
    newRow.getCell(1, 5).setValue(user.isAdmin.toLowerCase() == 'true');
    
    function setupSheet(sheetName, user) {
      var sheet = ss.getSheetByName(sheetName);
      var newRow = sheet.getRange(sheet.getLastRow()+1, 1, 1, 3);
      newRow.getCell(1, 1).setFormula("row()-2");
      newRow.getCell(1, 2).setValue(user.email);
      newRow.getCell(1, 3).setValue(user.name);
      userRow = sheet.getLastRow()-2;
    }
    if (!(user.isTeacher==="false" && user.isAdmin==='true')) { 
      setupSheet((user.isTeacher==="true"?"teacher": "student")+" data current",user); 
      setupSheet((user.isTeacher==="true"?"teacher": "student")+" data new",user);
    }
    newRow.getCell(1, 6).setValue(userRow);
  }
}

function editUser(user) {
  if (isAdmin()){
    openSpreadSheet();
    var sheet = ss.getSheetByName("Permission data");
    var userRow = sheet.getRange(user.id+1, 1, 1, 6);
    var curRow = user.userRow;
    
    function addToSheet(sheetName, user) {
      var sheet = ss.getSheetByName(sheetName);
      var newRow = sheet.getRange(sheet.getLastRow()+1, 1, 1, 3);
      newRow.getCell(1, 1).setValue(sheet.getLastRow()-1);
      newRow.getCell(1, 2).setValue(user.email);
      newRow.getCell(1, 3).setValue(user.name);
      curRow = sheet.getLastRow() - 2;
    }
    
    function removeFromSheet(sheetName) {
      var sheet = ss.getSheetByName(sheetName);
      sheet.deleteRow(curRow+2);
    }
    
    // enable when can use query to automaticly fix row numbers
    /*
    if (userRow.getCell(1, 4).getValue() != user.isTeacher) {
      if ((!userRow.getCell(1, 4).getValue() && !userRow.getCell(1, 5).getValue()) ||  userRow.getCell(1, 4).getValue()) {
        removeFromSheet((user.isTeacher === false ?"teacher": "student")+" data current");
        removeFromSheet((user.isTeacher === false ?"teacher": "student")+" data new");
      }
      addToSheet((user.isTeacher=== true ?"teacher": "student")+" data current",user); 
      addToSheet((user.isTeacher=== true ?"teacher": "student")+" data new",user);
      userRow.getCell(1, 6).setValue(curRow);
    }
    */
    userRow.getCell(1, 4).setValue(user.isTeacher === true);
    userRow.getCell(1, 5).setValue(user.isAdmin === true);
  }
  
}

function include(filename) {
  return HtmlService.createHtmlOutputFromFile(filename)
      .getContent();
}

function clearCache() {
  CacheService.getUserCache().remove('isAdmin');
  CacheService.getUserCache().remove('isTeacher');
  CacheService.getUserCache().remove('curUserRow');
}

function addTrigger() {
  ScriptApp.newTrigger('onEdit')
   .forSpreadsheet(sheetId)
   .onEdit()
   .create();
}