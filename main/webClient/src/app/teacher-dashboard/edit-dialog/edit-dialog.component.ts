import { Component, OnInit, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material';

import { TeacherService } from '../../teacher.service';
import { Block, Teacher } from '../../Interfaces';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-edit-dialog',
  templateUrl: './edit-dialog.component.html',
  styleUrls: ['./edit-dialog.component.css']
})
export class EditDialogComponent{

  block$: Observable<Block> = this.teacherService.getNextBlock(this.data.blockID);
  
  constructor(
    private teacherService: TeacherService,
    @Inject(MAT_DIALOG_DATA) public data: {
      blockID: number}
  ) { }
}