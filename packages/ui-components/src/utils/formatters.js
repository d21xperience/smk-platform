/**
 * Shared formatters untuk data display
 */

export const formatDate = (date, format = "DD MMMM YYYY") => {
  if (!date) return "-";

  const d = new Date(date);
  if (isNaN(d.getTime())) return "-";

  const months = [
    "Januari",
    "Februari",
    "Maret",
    "April",
    "Mei",
    "Juni",
    "Juli",
    "Agustus",
    "September",
    "Oktober",
    "November",
    "Desember",
  ];

  const day = d.getDate().toString().padStart(2, "0");
  const month = months[d.getMonth()];
  const year = d.getFullYear();

  switch (format) {
    case "DD MMMM YYYY":
      return `${day} ${month} ${year}`;
    case "DD/MM/YYYY":
      return `${day}/${(d.getMonth() + 1).toString().padStart(2, "0")}/${year}`;
    case "YYYY-MM-DD":
      return `${year}-${(d.getMonth() + 1).toString().padStart(2, "0")}-${day}`;
    default:
      return `${day} ${month} ${year}`;
  }
};

export const formatCurrency = (amount, currency = "Rp") => {
  if (amount === null || amount === undefined) return "-";

  const num = Number(amount);
  if (isNaN(num)) return "-";

  return `${currency} ${num.toLocaleString("id-ID")}`;
};

export const formatNumber = (num) => {
  if (num === null || num === undefined) return "-";

  const n = Number(num);
  if (isNaN(n)) return "-";

  return n.toLocaleString("id-ID");
};

export const truncateText = (text, maxLength = 100) => {
  if (!text) return "";
  if (text.length <= maxLength) return text;
  return text.substring(0, maxLength) + "...";
};

export const capitalize = (str) => {
  if (!str) return "";
  return str.charAt(0).toUpperCase() + str.slice(1).toLowerCase();
};

export const titleCase = (str) => {
  if (!str) return "";
  return str.toLowerCase().replace(/\b\w/g, (c) => c.toUpperCase());
};
