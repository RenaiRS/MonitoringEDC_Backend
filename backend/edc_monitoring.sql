-- phpMyAdmin SQL Dump
-- version 5.2.0
-- https://www.phpmyadmin.net/
--
-- Host: localhost:3306
-- Generation Time: Jan 08, 2026 at 08:11 AM
-- Server version: 8.0.30
-- PHP Version: 8.1.10

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `edc_monitoring`
--

-- --------------------------------------------------------

--
-- Table structure for table `mesin_edcs`
--

CREATE TABLE `mesin_edcs` (
  `id` bigint UNSIGNED NOT NULL,
  `terminal_id` varchar(191) NOT NULL,
  `m_id` longtext,
  `sn` longtext,
  `nama_nasabah` longtext,
  `kota` longtext,
  `cabang` longtext,
  `tipe_edc` longtext,
  `vendor` longtext,
  `status_mesin` longtext,
  `status_data` longtext,
  `lokasi_mesin` longtext,
  `letak_mesin` longtext,
  `tanggal_pasang` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `mesin_edcs`
--

INSERT INTO `mesin_edcs` (`id`, `terminal_id`, `m_id`, `sn`, `nama_nasabah`, `kota`, `cabang`, `tipe_edc`, `vendor`, `status_mesin`, `status_data`, `lokasi_mesin`, `letak_mesin`, `tanggal_pasang`, `created_at`, `updated_at`) VALUES
(1, '12345678', 'MID001122', '', '', 'Palembang', 'Cabang Veteran', 'Verifone VX520', '', 'aktif', 'vendor_only', '', 'nasabah', '2024-06-01 07:00:00.000', '2026-01-07 09:25:02.617', '2026-01-07 09:25:02.617');

-- --------------------------------------------------------

--
-- Table structure for table `overdues`
--

CREATE TABLE `overdues` (
  `id` bigint UNSIGNED NOT NULL,
  `mesin_id` bigint UNSIGNED DEFAULT NULL,
  `status_perbaikan` longtext,
  `estimasi_selesai` datetime(3) DEFAULT NULL,
  `terlambat_hari` bigint DEFAULT NULL,
  `estimasi_kerugian` bigint DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `perbaikans`
--

CREATE TABLE `perbaikans` (
  `id` bigint UNSIGNED NOT NULL,
  `mesin_id` bigint UNSIGNED DEFAULT NULL,
  `status_perbaikan` longtext,
  `estimasi_perbaikan` datetime(3) DEFAULT NULL,
  `estimasi_selesai` datetime(3) DEFAULT NULL,
  `sn` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `sewas`
--

CREATE TABLE `sewas` (
  `id` bigint UNSIGNED NOT NULL,
  `mesin_id` bigint UNSIGNED DEFAULT NULL,
  `status_sewa` longtext,
  `biaya_bulanan` bigint DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `sewas`
--

INSERT INTO `sewas` (`id`, `mesin_id`, `status_sewa`, `biaya_bulanan`, `created_at`, `updated_at`) VALUES
(1, 1, 'berakhir', 1500000, '2026-01-07 09:25:02.623', '2026-01-07 09:25:02.623');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `mesin_edcs`
--
ALTER TABLE `mesin_edcs`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uni_mesin_edcs_terminal_id` (`terminal_id`);

--
-- Indexes for table `overdues`
--
ALTER TABLE `overdues`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_overdues_mesin` (`mesin_id`);

--
-- Indexes for table `perbaikans`
--
ALTER TABLE `perbaikans`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_mesin_edcs_perbaikan` (`mesin_id`);

--
-- Indexes for table `sewas`
--
ALTER TABLE `sewas`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_mesin_edcs_sewa` (`mesin_id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `mesin_edcs`
--
ALTER TABLE `mesin_edcs`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `overdues`
--
ALTER TABLE `overdues`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `perbaikans`
--
ALTER TABLE `perbaikans`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `sewas`
--
ALTER TABLE `sewas`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `overdues`
--
ALTER TABLE `overdues`
  ADD CONSTRAINT `fk_overdues_mesin` FOREIGN KEY (`mesin_id`) REFERENCES `mesin_edcs` (`id`);

--
-- Constraints for table `perbaikans`
--
ALTER TABLE `perbaikans`
  ADD CONSTRAINT `fk_mesin_edcs_perbaikan` FOREIGN KEY (`mesin_id`) REFERENCES `mesin_edcs` (`id`);

--
-- Constraints for table `sewas`
--
ALTER TABLE `sewas`
  ADD CONSTRAINT `fk_mesin_edcs_sewa` FOREIGN KEY (`mesin_id`) REFERENCES `mesin_edcs` (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
