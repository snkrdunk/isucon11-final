INSERT INTO `users` (`id`, `name`, `mail_address`, `hashed_password`, `type`) VALUES
(
  '01234567-89ab-cdef-0001-000000000001',
  '佐藤太郎',
  'sato@example.com',
  '$2a$10$KEgha.chGu1/N4kHZ./rIeK1QISkv8sYk15Mqktr6BGB8xomRRe02', -- "password"
  'student'
),
(
  '01234567-89ab-cdef-0001-000000000002',
  '鈴木次郎',
  'suzuki@example.com',
  '$2a$10$KEgha.chGu1/N4kHZ./rIeK1QISkv8sYk15Mqktr6BGB8xomRRe02',
  'student'
),
(
  '01234567-89ab-cdef-0001-000000000003',
  '高橋三郎',
  'takahashi@example.com',
  '$2a$10$KEgha.chGu1/N4kHZ./rIeK1QISkv8sYk15Mqktr6BGB8xomRRe02',
  'student'
),
(
  '01234567-89ab-cdef-0001-000000000004',
  '椅子昆',
  'isu@example.com',
  '$2a$10$KEgha.chGu1/N4kHZ./rIeK1QISkv8sYk15Mqktr6BGB8xomRRe02',
  'faculty'
),
(
  '2a931c3f-0576-464a-a05d-6e700ad54d70',
  '田山勝蔵',
  'tayama.s@isu.com',
  '$2a$10$KEgha.chGu1/N4kHZ./rIeK1QISkv8sYk15Mqktr6BGB8xomRRe02',
  'faculty'
);

INSERT INTO `courses` (`id`, `code`, `type`, `name`, `description`, `credit`, `classroom`, `capacity`, `teacher_id`, `keywords`) VALUES
(
  '01234567-89ab-cdef-0002-000000000001',
  'ISU.F117',
  'liberal-arts',
  '微分積分基礎',
  '微積分の基礎を学びます。',
  2,
  'A101講義室',
  100,
  '01234567-89ab-cdef-0001-000000000004',
  '数学 微分 積分 基礎'
),
(
  '01234567-89ab-cdef-0002-000000000002',
  'ISU.M101',
  'liberal-arts',
  '線形代数基礎',
  '線形代数の基礎を学びます。',
  2,
  'B101講義室',
  100,
  '2a931c3f-0576-464a-a05d-6e700ad54d70',
  '数学 線形代数 基礎'
),
(
  '01234567-89ab-cdef-0002-000000000003',
  'CSC.A331',
  'major-subjects',
  'アルゴリズム基礎',
  'アルゴリズムの基礎を学びます。',
  2,
  'C101講義室',
  NULL,
  '2a931c3f-0576-464a-a05d-6e700ad54d70',
  '計算機科学 アルゴリズム 基礎'
),
(
  '01234567-89ab-cdef-0002-000000000004',
  'CSC.A332',
  'major-subjects',
  'アルゴリズム応用',
  'アルゴリズムの応用を学びます。',
  2,
  'C101講義室',
  30,
  '2a931c3f-0576-464a-a05d-6e700ad54d70',
  '計算機科学 アルゴリズム 応用'
),
(
  '01234567-89ab-cdef-0002-000000000011',
  'ISU.F118',
  'major-subjects',
  '微分積分応用',
  '微積分の応用を学びます。',
  2,
  'A102講義室',
  100,
  '01234567-89ab-cdef-0001-000000000004',
  '数学 微分 積分 応用'
),
(
  '01234567-89ab-cdef-0002-000000000012',
  'ISU.F119',
  'major-subjects',
  '線形代数応用',
  '線形代数の応用を学びます。',
  2,
  'B102講義室',
  NULL,
  '01234567-89ab-cdef-0001-000000000004',
  '数学 線形代数 応用'
),
(
  '01234567-89ab-cdef-0002-000000000013',
  'CSP.B003',
  'liberal-arts',
  'プログラミング',
  'プログラミングを学びます。',
  2,
  'C102講義室',
  100,
  '2a931c3f-0576-464a-a05d-6e700ad54d70',
  '計算機科学 C言語'
),
(
  '01234567-89ab-cdef-0002-000000000014',
  'CSC.B103',
  'major-subjects',
  'プログラミング演習A',
  'プログラミングの演習を行います。',
  1,
  '演習室1',
  50,
  '2a931c3f-0576-464a-a05d-6e700ad54d70',
  '計算機科学 C言語 演習'
),
(
  '01234567-89ab-cdef-0002-000000000015',
  'CSC.B104',
  'major-subjects',
  'プログラミング演習B',
  'プログラミングの演習を行います。',
  1,
  '演習室1',
  50,
  '2a931c3f-0576-464a-a05d-6e700ad54d70',
  '計算機科学 C言語 演習'
);

INSERT INTO `course_requirements` (`course_id`, `required_course_id`) VALUES
(
  '01234567-89ab-cdef-0002-000000000011',
  '01234567-89ab-cdef-0002-000000000001'
),
(
  '01234567-89ab-cdef-0002-000000000012',
  '01234567-89ab-cdef-0002-000000000002'
);

INSERT INTO `schedules` (`id`, `period`, `day_of_week`, `semester`, `year`) VALUES
(
  '01234567-89ab-cdef-0003-000000000001',
  1,
  'monday',
  'first',
  2021
),
(
  '01234567-89ab-cdef-0003-000000000002',
  2,
  'wednesday',
  'first',
  2021
),
(
  '01234567-89ab-cdef-0003-000000000003',
  3,
  'friday',
  'first',
  2021
),
(
  '01234567-89ab-cdef-0003-000000000011',
  1,
  'tuesday',
  'second',
  2021
),
(
  '01234567-89ab-cdef-0003-000000000012',
  2,
  'thursday',
  'second',
  2021
),
(
  '01234567-89ab-cdef-0003-000000000013',
  3,
  'friday',
  'second',
  2021
),
(
  '01234567-89ab-cdef-0003-000000000014',
  4,
  'friday',
  'second',
  2021
),
(
  '01234567-89ab-cdef-0003-000000000015',
  4,
  'friday',
  'second',
  2021
);

INSERT INTO `course_schedules` (`course_id`, `schedule_id`) VALUES
(
  '01234567-89ab-cdef-0002-000000000001',
  '01234567-89ab-cdef-0003-000000000001'
),
(
  '01234567-89ab-cdef-0002-000000000002',
  '01234567-89ab-cdef-0003-000000000002'
),
(
  '01234567-89ab-cdef-0002-000000000003',
  '01234567-89ab-cdef-0003-000000000003'
),
(
  '01234567-89ab-cdef-0002-000000000011',
  '01234567-89ab-cdef-0003-000000000011'
),
(
  '01234567-89ab-cdef-0002-000000000012',
  '01234567-89ab-cdef-0003-000000000012'
),
(
  '01234567-89ab-cdef-0002-000000000013',
  '01234567-89ab-cdef-0003-000000000013'
),
(
  '01234567-89ab-cdef-0002-000000000014',
  '01234567-89ab-cdef-0003-000000000014'
),
(
  '01234567-89ab-cdef-0002-000000000015',
  '01234567-89ab-cdef-0003-000000000015'
);
