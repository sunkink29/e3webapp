app.service('objService', function () {
  this.getTeacherObjFromList = function(teacherList) {
    var teacherObj = {name:"Unassigned", id: -1, blocks: [{status: 'open'}, {status: 'open'}]};
    if (teacherList != null) {
      teacherObj = {
        id: teacherList[0],
        email: teacherList[1],
        name: teacherList[2],
        blocks: [this.getBlockInfoObjFromList(teacherList.slice(3,9)),
        this.getBlockInfoObjFromList(teacherList.slice(9,15))]
      }
    }
    return teacherObj;
  }
  this.getStudentObjFromList = function(studentList) {
    var studentObj = {
      id: studentList[0],
      email: studentList[1],
      name: studentList[2],
      grade: studentList[3],
      blocks: [this.getTeacherObjFromList(studentList[4]), 
      this.getTeacherObjFromList(studentList[5])]
    }
    return studentObj;
  }
  this.getBlockInfoObjFromList = function(blockInfoList) {
    var blockInfoObj = {
      subject: blockInfoList[0],
      description: blockInfoList[1],
      roomNumber: blockInfoList[2],
      curSize: blockInfoList[3],
      maxSize: blockInfoList[4],
      status: blockInfoList[5]
    }
    return blockInfoObj;
  }
  this.getUserObjFromList = function(userList) {
    var userObj = {
      id: userList[0],
      userRow: userList[5],
      email: userList[1],
      name: userList[2],
      isTeacher: userList[3],
      isAdmin: userList[4]
    }
    return userObj;
  }
});