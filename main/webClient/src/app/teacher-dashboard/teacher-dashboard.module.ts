import { NgModule } from '@angular/core';
import { CommonModule, NgForOf } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { TeacherDashboardRoutingModule } from './teacher-dashboard-routing.module';
import { TeacherDashboardComponent } from './teacher-dashboard/teacher-dashboard.component';
import { MaterialImportModule } from '../material-import.module';
import { EditDialogComponent } from './edit-dialog/edit-dialog.component';
import { AddDialogComponent } from './add-dialog/add-dialog.component';

@NgModule({
  imports: [
    CommonModule,
    TeacherDashboardRoutingModule,
    MaterialImportModule,
    FormsModule
  ],
  entryComponents: [
    EditDialogComponent,
    AddDialogComponent
  ],
  declarations: [TeacherDashboardComponent, EditDialogComponent, AddDialogComponent]
})
export class TeacherDashboardModule { }
