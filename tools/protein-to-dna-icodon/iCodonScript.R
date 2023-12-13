args <- commandArgs(trailingOnly = TRUE)

library(iCodon)

sequence <- readLines("/outputs/dna_sequence_after_reverse_translation.txt", warn=FALSE)

specie <- args[1]
iterations <- as.integer(args[2])
make_more_optimal <- as.logical(args[3])

result <- run_optimization_shinny(sequence, specie)
output_file <- "/outputs/optimized_shiny.csv"
write.csv(result, output_file, row.names = FALSE)

print(paste("Optimizing for species:", specie, "for", iterations, "iterations, with make_more_optimal:", ifelse(make_more_optimal, "T", "F")))

result <- optimizer(sequence, specie, iterations, make_more_optimal)
output_file <- "/outputs/optimizer_result.csv"
write.csv(result, output_file, row.names = FALSE)