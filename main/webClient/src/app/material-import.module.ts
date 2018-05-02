import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import {MatToolbarModule, MatTableModule, MatButtonModule, MatDialogModule} from '@angular/material';

@NgModule({
  exports: [
    MatToolbarModule,
    MatTableModule,
    MatButtonModule,
    MatDialogModule
    ]
})
export class MaterialImportModule {}
