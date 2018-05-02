import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { StudentDashboardComponent } from './student-dashboard/student-dashboard.component';

const routes: Routes = [
  {
    path: '',
    component: StudentDashboardComponent
  },
  {
    path: 'change/:id',
    component: StudentDashboardComponent
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class StudentDashboardRoutingModule { }
