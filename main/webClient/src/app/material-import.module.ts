import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import {MatToolbarModule, MatTableModule, MatButtonModule, MatDialogModule, MatTabsModule, MatFormFieldModule, MatInputModule, MatCheckboxModule} from '@angular/material';

@NgModule({
  exports: [
    MatToolbarModule,
    MatTableModule,
    MatButtonModule,
    MatDialogModule,
    MatTabsModule,
    MatFormFieldModule,
    MatInputModule,
    MatCheckboxModule
    ]
})
export class MaterialImportModule {}
