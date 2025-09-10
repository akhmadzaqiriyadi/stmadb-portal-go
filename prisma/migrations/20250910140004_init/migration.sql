-- CreateTable
CREATE TABLE `users` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(100) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `role` ENUM('admin', 'teacher', 'student', 'staff') NOT NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT true,
    `last_login` DATETIME(3) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    UNIQUE INDEX `users_username_key`(`username`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `teachers` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT NOT NULL,
    `nip` VARCHAR(18) NULL,
    `nik` VARCHAR(16) NULL,
    `full_name` VARCHAR(255) NOT NULL,
    `phone_number` VARCHAR(20) NULL,
    `employment_status` ENUM('ASN', 'GTT', 'PTT', 'Tetap') NOT NULL,
    `signature_image_path` VARCHAR(255) NULL,

    UNIQUE INDEX `teachers_user_id_key`(`user_id`),
    UNIQUE INDEX `teachers_nip_key`(`nip`),
    UNIQUE INDEX `teachers_nik_key`(`nik`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `subjects` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `subject_code` VARCHAR(20) NOT NULL,
    `subject_name` VARCHAR(255) NOT NULL,

    UNIQUE INDEX `subjects_subject_code_key`(`subject_code`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `classes` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `class_name` VARCHAR(100) NOT NULL,
    `grade_level` VARCHAR(10) NOT NULL,
    `major` VARCHAR(100) NULL,
    `academic_year` VARCHAR(10) NOT NULL,
    `homeroom_teacher_id` BIGINT NULL,
    `counselor_id` BIGINT NULL,

    INDEX `classes_homeroom_teacher_id_idx`(`homeroom_teacher_id`),
    INDEX `classes_counselor_id_idx`(`counselor_id`),
    UNIQUE INDEX `classes_class_name_academic_year_key`(`class_name`, `academic_year`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `students` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT NOT NULL,
    `current_class_id` BIGINT NULL,
    `nis` VARCHAR(16) NOT NULL,
    `nisn` VARCHAR(10) NULL,
    `full_name` VARCHAR(255) NOT NULL,
    `gender` ENUM('L', 'P') NOT NULL,
    `address` TEXT NULL,
    `phone_number` VARCHAR(20) NULL,
    `status` ENUM('AKTIF', 'LULUS', 'PINDAH', 'DO') NOT NULL DEFAULT 'AKTIF',
    `rfid_uid` VARCHAR(100) NULL,

    UNIQUE INDEX `students_user_id_key`(`user_id`),
    UNIQUE INDEX `students_nis_key`(`nis`),
    UNIQUE INDEX `students_nisn_key`(`nisn`),
    UNIQUE INDEX `students_rfid_uid_key`(`rfid_uid`),
    INDEX `students_current_class_id_idx`(`current_class_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `schedules` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `class_id` BIGINT NOT NULL,
    `subject_id` BIGINT NOT NULL,
    `teacher_id` BIGINT NOT NULL,
    `day_of_week` ENUM('Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu', 'Minggu') NOT NULL,
    `start_time` TIME NOT NULL,
    `end_time` TIME NOT NULL,
    `room` VARCHAR(50) NULL,

    INDEX `schedules_class_id_idx`(`class_id`),
    INDEX `schedules_subject_id_idx`(`subject_id`),
    INDEX `schedules_teacher_id_idx`(`teacher_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `teaching_journals` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `schedule_id` BIGINT NOT NULL,
    `teaching_date` DATE NOT NULL,
    `topic` TEXT NOT NULL,
    `student_attendance_summary` VARCHAR(255) NULL,
    `notes` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `teaching_journals_schedule_id_idx`(`schedule_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `companies` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL,
    `address` TEXT NULL,
    `coordinates` VARCHAR(100) NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `internship_placements` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `student_id` BIGINT NOT NULL,
    `company_id` BIGINT NOT NULL,
    `supervisor_teacher_id` BIGINT NULL,
    `start_date` DATE NOT NULL,
    `end_date` DATE NULL,
    `status` ENUM('Aktif', 'Selesai', 'Batal') NOT NULL DEFAULT 'Aktif',

    INDEX `internship_placements_student_id_idx`(`student_id`),
    INDEX `internship_placements_company_id_idx`(`company_id`),
    INDEX `internship_placements_supervisor_teacher_id_idx`(`supervisor_teacher_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `internship_journals` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `placement_id` BIGINT NOT NULL,
    `activity_date` DATE NOT NULL,
    `activity_description` TEXT NOT NULL,
    `status` ENUM('Pending', 'Approved', 'Rejected') NOT NULL DEFAULT 'Pending',
    `supervisor_notes` TEXT NULL,
    `approved_at` DATETIME(3) NULL,

    INDEX `internship_journals_placement_id_idx`(`placement_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `attendances` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT NOT NULL,
    `timestamp` DATETIME(3) NOT NULL,
    `status` ENUM('Masuk', 'Pulang') NOT NULL,
    `location_coordinates` VARCHAR(100) NULL,
    `photo_path` VARCHAR(255) NULL,

    INDEX `attendances_user_id_timestamp_idx`(`user_id`, `timestamp`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `leave_requests` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT NOT NULL,
    `request_type` ENUM('Sakit', 'Izin', 'Cuti', 'DinasLuar') NOT NULL,
    `start_date` DATE NOT NULL,
    `end_date` DATE NOT NULL,
    `reason` TEXT NOT NULL,
    `attachment_path` VARCHAR(255) NULL,
    `status` ENUM('Pending', 'Approved', 'Rejected') NOT NULL DEFAULT 'Pending',
    `verifier_id` BIGINT NULL,
    `verified_at` DATETIME(3) NULL,
    `rejection_reason` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `leave_requests_user_id_idx`(`user_id`),
    INDEX `leave_requests_verifier_id_idx`(`verifier_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ramadan_activities` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `student_id` BIGINT NOT NULL,
    `activity_date` DATE NOT NULL,
    `activity_type` ENUM('Puasa', 'SalatFardu', 'SalatTarawih', 'SalatSunnah', 'Tadarus', 'Taklim', 'SalatJumat') NOT NULL,
    `sub_type` VARCHAR(50) NULL,
    `is_done` BOOLEAN NOT NULL DEFAULT false,
    `description` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `ramadan_activities_student_id_idx`(`student_id`),
    UNIQUE INDEX `ramadan_activities_student_id_activity_date_activity_type_su_key`(`student_id`, `activity_date`, `activity_type`, `sub_type`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `queue_counters` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `counter_name` VARCHAR(255) NOT NULL,
    `counter_code` VARCHAR(10) NOT NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT true,

    UNIQUE INDEX `queue_counters_counter_code_key`(`counter_code`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `queue_tickets` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `counter_id` INTEGER NOT NULL,
    `student_id` BIGINT NULL,
    `ticket_number` INTEGER NOT NULL,
    `queue_date` DATE NOT NULL,
    `status` ENUM('Waiting', 'Called', 'Serving', 'Finished', 'Skipped') NOT NULL DEFAULT 'Waiting',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `called_at` DATETIME(3) NULL,
    `finished_at` DATETIME(3) NULL,

    INDEX `queue_tickets_student_id_idx`(`student_id`),
    UNIQUE INDEX `queue_tickets_queue_date_counter_id_ticket_number_key`(`queue_date`, `counter_id`, `ticket_number`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `exams` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `exam_name` VARCHAR(255) NOT NULL,
    `start_date` DATE NOT NULL,
    `end_date` DATE NOT NULL,
    `academic_year` VARCHAR(10) NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `exam_schedules` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `exam_id` BIGINT NOT NULL,
    `subject_id` BIGINT NOT NULL,
    `exam_date` DATE NOT NULL,
    `start_time` TIME NOT NULL,
    `end_time` TIME NOT NULL,
    `session` INTEGER NOT NULL DEFAULT 1,

    INDEX `exam_schedules_exam_id_idx`(`exam_id`),
    INDEX `exam_schedules_subject_id_idx`(`subject_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `exam_rooms` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `room_name` VARCHAR(100) NOT NULL,
    `capacity` INTEGER NOT NULL,

    UNIQUE INDEX `exam_rooms_room_name_key`(`room_name`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `exam_assignments` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `schedule_id` BIGINT NOT NULL,
    `student_id` BIGINT NOT NULL,
    `room_id` INTEGER NOT NULL,
    `supervisor_1_id` BIGINT NULL,
    `supervisor_2_id` BIGINT NULL,
    `token` VARCHAR(10) NULL,
    `student_attendance` ENUM('Hadir', 'Absen', 'Sakit', 'Izin') NULL,

    UNIQUE INDEX `exam_assignments_token_key`(`token`),
    INDEX `exam_assignments_schedule_id_idx`(`schedule_id`),
    INDEX `exam_assignments_student_id_idx`(`student_id`),
    INDEX `exam_assignments_room_id_idx`(`room_id`),
    INDEX `exam_assignments_supervisor_1_id_idx`(`supervisor_1_id`),
    INDEX `exam_assignments_supervisor_2_id_idx`(`supervisor_2_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `exam_incident_reports` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `exam_id` BIGINT NOT NULL,
    `room_id` INTEGER NOT NULL,
    `session` INTEGER NOT NULL,
    `report_date` DATE NOT NULL,
    `notes` TEXT NULL,
    `absent_students_summary` TEXT NULL,
    `reported_by_id` BIGINT NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `exam_incident_reports_exam_id_idx`(`exam_id`),
    INDEX `exam_incident_reports_room_id_idx`(`room_id`),
    INDEX `exam_incident_reports_reported_by_id_idx`(`reported_by_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- AddForeignKey
ALTER TABLE `teachers` ADD CONSTRAINT `teachers_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `classes` ADD CONSTRAINT `classes_homeroom_teacher_id_fkey` FOREIGN KEY (`homeroom_teacher_id`) REFERENCES `teachers`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `classes` ADD CONSTRAINT `classes_counselor_id_fkey` FOREIGN KEY (`counselor_id`) REFERENCES `teachers`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `students` ADD CONSTRAINT `students_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `students` ADD CONSTRAINT `students_current_class_id_fkey` FOREIGN KEY (`current_class_id`) REFERENCES `classes`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `schedules` ADD CONSTRAINT `schedules_class_id_fkey` FOREIGN KEY (`class_id`) REFERENCES `classes`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `schedules` ADD CONSTRAINT `schedules_subject_id_fkey` FOREIGN KEY (`subject_id`) REFERENCES `subjects`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `schedules` ADD CONSTRAINT `schedules_teacher_id_fkey` FOREIGN KEY (`teacher_id`) REFERENCES `teachers`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `teaching_journals` ADD CONSTRAINT `teaching_journals_schedule_id_fkey` FOREIGN KEY (`schedule_id`) REFERENCES `schedules`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `internship_placements` ADD CONSTRAINT `internship_placements_student_id_fkey` FOREIGN KEY (`student_id`) REFERENCES `students`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `internship_placements` ADD CONSTRAINT `internship_placements_company_id_fkey` FOREIGN KEY (`company_id`) REFERENCES `companies`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `internship_placements` ADD CONSTRAINT `internship_placements_supervisor_teacher_id_fkey` FOREIGN KEY (`supervisor_teacher_id`) REFERENCES `teachers`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `internship_journals` ADD CONSTRAINT `internship_journals_placement_id_fkey` FOREIGN KEY (`placement_id`) REFERENCES `internship_placements`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `attendances` ADD CONSTRAINT `attendances_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `leave_requests` ADD CONSTRAINT `leave_requests_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `leave_requests` ADD CONSTRAINT `leave_requests_verifier_id_fkey` FOREIGN KEY (`verifier_id`) REFERENCES `users`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ramadan_activities` ADD CONSTRAINT `ramadan_activities_student_id_fkey` FOREIGN KEY (`student_id`) REFERENCES `students`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `queue_tickets` ADD CONSTRAINT `queue_tickets_counter_id_fkey` FOREIGN KEY (`counter_id`) REFERENCES `queue_counters`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `queue_tickets` ADD CONSTRAINT `queue_tickets_student_id_fkey` FOREIGN KEY (`student_id`) REFERENCES `students`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `exam_schedules` ADD CONSTRAINT `exam_schedules_exam_id_fkey` FOREIGN KEY (`exam_id`) REFERENCES `exams`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `exam_schedules` ADD CONSTRAINT `exam_schedules_subject_id_fkey` FOREIGN KEY (`subject_id`) REFERENCES `subjects`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `exam_assignments` ADD CONSTRAINT `exam_assignments_schedule_id_fkey` FOREIGN KEY (`schedule_id`) REFERENCES `exam_schedules`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `exam_assignments` ADD CONSTRAINT `exam_assignments_student_id_fkey` FOREIGN KEY (`student_id`) REFERENCES `students`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `exam_assignments` ADD CONSTRAINT `exam_assignments_room_id_fkey` FOREIGN KEY (`room_id`) REFERENCES `exam_rooms`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `exam_assignments` ADD CONSTRAINT `exam_assignments_supervisor_1_id_fkey` FOREIGN KEY (`supervisor_1_id`) REFERENCES `teachers`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `exam_assignments` ADD CONSTRAINT `exam_assignments_supervisor_2_id_fkey` FOREIGN KEY (`supervisor_2_id`) REFERENCES `teachers`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `exam_incident_reports` ADD CONSTRAINT `exam_incident_reports_exam_id_fkey` FOREIGN KEY (`exam_id`) REFERENCES `exams`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `exam_incident_reports` ADD CONSTRAINT `exam_incident_reports_room_id_fkey` FOREIGN KEY (`room_id`) REFERENCES `exam_rooms`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `exam_incident_reports` ADD CONSTRAINT `exam_incident_reports_reported_by_id_fkey` FOREIGN KEY (`reported_by_id`) REFERENCES `teachers`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;
