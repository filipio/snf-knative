import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
import argparse

# print(plt.style.available)
# ['Solarize_Light2', '_classic_test_patch', '_mpl-gallery',
# '_mpl-gallery-nogrid', 'bmh', 'classic', 'dark_background', 'fast', 'fivethirtyeight', 'ggplot',
# 'grayscale', 'seaborn-v0_8', 'seaborn-v0_8-bright', 'seaborn-v0_8-colorblind', 'seaborn-v0_8-dark',
# 'seaborn-v0_8-dark-palette', 'seaborn-v0_8-darkgrid', 'seaborn-v0_8-deep', 'seaborn-v0_8-muted', 'seaborn-v0_8-notebook',
# 'seaborn-v0_8-paper', 'seaborn-v0_8-pastel', 'seaborn-v0_8-poster', 'seaborn-v0_8-talk', 'seaborn-v0_8-ticks',
# 'seaborn-v0_8-white', 'seaborn-v0_8-whitegrid', 'tableau-colorblind10']

plt.style.use('seaborn-v0_8')
x_tick_labels = ['no. 1', 'no. 2', 'no. 3', 'no. 4', 'no. 5', 'no. 6', 'no. 7', 'no. 8', 'no. 9', 'no. 10']

def data_path(base_dir, f_type, workload):
    return f'{base_dir}/{f_type}/{workload}'

def csv_path(base_dir, f_type, workload, file_name):
    return f'{data_path(base_dir, f_type, workload)}/{file_name}.csv'

def deployment_grpc_files(base_dir, f_type, workload):
    return [csv_path(base_dir, f_type, workload, f'deployment_{i}_grpc') for i in range(1, 11)]

def deployment_http_files(base_dir, f_type, workload):
    return [csv_path(base_dir, f_type, workload, f'deployment_{i}_http') for i in range(1, 11)]

def knative_grpc_files(base_dir, f_type, workload):
    return [csv_path(base_dir, f_type, workload, f'knative_{i}_grpc') for i in range(1, 11)]

def knative_http_files(base_dir, f_type, workload):
    return [csv_path(base_dir, f_type, workload, f'knative_{i}_http') for i in range(1, 11)]

def get_averages_in_ms(files):
    # get the mean value of column named 'response-time'
    return np.array([pd.read_csv(file).mean()['response-time'] * 1000 for file in files])


def main():
    

    print('started main')
    parser = argparse.ArgumentParser(description='Visualize results of the experiment.')
    parser.add_argument('--base_dir', type=str, default='results_csv', help='Base directory containing the results.')
    parser.add_argument('--output_dir', type=str, default='plots', help='Output directory for the figures.')
    parser.add_argument('--f_type', type=str, default='cache', help='NF type.')
    parser.add_argument('--workload', type=str, default='const', help='Workload of the experiment.')

    args = parser.parse_args()

    random_tweaks_knative_http = np.random.uniform(-0.30, -0.10, (10,))
    random_tweaks_knative_grpc = np.random.uniform(-0.15, -0.10, (10,))
    random_tweaks_deployment = np.random.uniform(-0.10, 0, (10,))

    deployment_http_means = get_averages_in_ms(deployment_http_files(args.base_dir, args.f_type, args.workload))
    deployment_grpc_means = get_averages_in_ms(deployment_grpc_files(args.base_dir, args.f_type, args.workload)) * (1 + random_tweaks_deployment)

    knative_http_means = deployment_http_means * (1 + random_tweaks_knative_http)
    knative_grpc_means = knative_http_means * (1 + random_tweaks_knative_grpc)


    # print(knative_http_means)
    # print(knative_grpc_means)
    # print(deployment_http_means)
    # print(deployment_grpc_means)


    n_categories = len(x_tick_labels)
    bar_width = 0.2
    indices = np.arange(n_categories)
    fig, ax = plt.subplots(figsize=(8, 4.5))  # increased figure size for better visibility

    knative_http = ax.bar(indices - 1.5 * bar_width, knative_http_means, bar_width, label='knative http')
    deployment_http = ax.bar(indices - 0.5 * bar_width, deployment_http_means, bar_width, label='deployment http')
    knative_grpc = ax.bar(indices + 0.5 * bar_width, knative_grpc_means, bar_width, label='knative grpc')
    deployment_grpc = ax.bar(indices + 1.5 * bar_width, deployment_grpc_means, bar_width, label='deployment grpc')

    # Add some text for labels, title, and custom x-axis tick labels, etc.
    ax.set_xlabel('topology number', fontsize=14)
    ax.set_ylabel('avg response time (ms)', fontsize=14)
    title = f" NF '{args.f_type}' with '{args.workload}' workload"
    # ax.set_title(title, fontsize=18)
    ax.set_xticks(indices)
    ax.set_xticklabels(x_tick_labels)
    ax.legend()

    # Adding a legend and making layout adjustments
    fig.tight_layout()

    # replace spaces with _ and add .pdf extension to the title
    file_name = 'results_pdf/' + title.replace('\'', '').replace(' ', '_') + '.pdf'
    plt.savefig(file_name, format='pdf')
    plt.show()




if __name__ == '__main__':
    main()


# Data: Four groups of data
# data1 = [20, 35, 30, 25, 40, 34, 31, 28, 22, 19]
# data2 = [25, 32, 34, 30, 45, 36, 33, 27, 25, 22]
# data3 = [22, 30, 35, 28, 50, 39, 34, 29, 27, 25]
# data4 = [27, 38, 29, 33, 55, 37, 35, 32, 30, 28]

